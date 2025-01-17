package analyzer

import (
	"errors"

	"github.com/yaklang/yaklang/common/sca/dxtypes"

	"github.com/yaklang/yaklang/common/sca/analyzer/dep-parser/golang/binary"
	"github.com/yaklang/yaklang/common/utils"
)

const (
	TypGoBinary TypAnalyzer = "go-binary-lang"

	statusExecutable int = 1
)

func init() {
	RegisterAnalyzer(TypGoBinary, NewGoBinaryAnalyzer())
}

type goBinaryAnalyzer struct{}

func NewGoBinaryAnalyzer() *goBinaryAnalyzer {
	return &goBinaryAnalyzer{}
}

func (a goBinaryAnalyzer) Match(info MatchInfo) int {
	if IsExecutable(info.FileHeader) {
		return statusExecutable
	}
	return 0
}

func (a goBinaryAnalyzer) Analyze(afi AnalyzeFileInfo) ([]*dxtypes.Package, error) {
	fi := afi.Self
	switch fi.MatchStatus {
	case statusExecutable:
		p := binary.NewParser()
		pkgs, err := ParseLanguageConfiguration(fi, p)
		if errors.Is(err, binary.ErrUnrecognizedExe) || errors.Is(err, binary.ErrNonGoBinary) {
			return nil, nil
		} else if err != nil {
			err = utils.Errorf("go binary parse error: %s", err)
		}
		return pkgs, err
	}

	return nil, nil
}
