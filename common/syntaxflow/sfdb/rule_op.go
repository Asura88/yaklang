package sfdb

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/schema"
	"github.com/yaklang/yaklang/common/syntaxflow/sfvm"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/bizhelper"
	"github.com/yaklang/yaklang/common/utils/filesys"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	"path"
	"strings"
)

func CreateOrUpdateSyntaxFlow(hash string, i any) error {
	db := consts.GetGormProfileDatabase()
	var rule schema.SyntaxFlowRule
	if err := db.Where("hash = ?", hash).First(&rule).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return db.Create(i).Error
		}
		return err
	}

	return db.Model(&rule).Updates(i).Error
}

func ImportValidRule(system filesys.FileSystem, ruleName string, content string) error {
	var language consts.Language
	languageRaw, _, _ := strings.Cut(ruleName, "-")
	switch strings.TrimSpace(strings.ToLower(languageRaw)) {
	case "yak", "yaklang":
		language = ssaapi.Yak
	case "java":
		language = ssaapi.JAVA
	case "php":
		language = ssaapi.PHP
	case "js", "es", "javascript", "ecmascript", "nodejs", "node", "node.js":
		language = ssaapi.JS
	}

	var ruleType schema.SyntaxFlowRuleType
	switch path.Ext(ruleName) {
	case ".sf", ".syntaxflow":
		ruleType = schema.SFR_RULE_TYPE_SF
	default:
		return utils.Errorf("invalid rule type: %v is not supported yet", ruleName)
	}

	frame, err := sfvm.NewSyntaxFlowVirtualMachine().Compile(content)
	if err != nil {
		return err
	}

	rule := &schema.SyntaxFlowRule{
		Language:    string(language),
		Title:       frame.Title,
		Description: frame.Description,
		Type:        ruleType,
		Content:     content,
		Purpose:     schema.ValidPurpose(frame.Purpose),
	}

	err = LoadFileSystem(rule, system)
	if err != nil {
		return utils.Wrap(err, "load file system error")
	}

	err = Valid(rule)
	if err != nil {
		return utils.Wrap(err, "valid rule error")
	}

	err = CreateOrUpdateSyntaxFlow(rule.CalcHash(), rule)
	if err != nil {
		return utils.Wrap(err, "create or update syntax flow rule error")
	}
	return nil
}

func YieldSyntaxFlowRules(db *gorm.DB, ctx context.Context) chan *schema.SyntaxFlowRule {
	outC := make(chan *schema.SyntaxFlowRule)
	go func() {
		defer close(outC)

		var page = 1
		for {
			var items []*schema.SyntaxFlowRule
			if _, b := bizhelper.Paging(db, page, 1000, &items); b.Error != nil {
				log.Errorf("paging failed: %s", b.Error)
				return
			}

			page++

			for _, d := range items {
				select {
				case <-ctx.Done():
					return
				case outC <- d:
				}
			}

			if len(items) < 1000 {
				return
			}
		}
	}()
	return outC
}