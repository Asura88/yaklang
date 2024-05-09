package javascript

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaklang/yaklang/common/yak/ssaapi"
	"testing"
)

func Test_JS_XMLHttpRequest(t *testing.T) {
	t.Run("simple get request", func(t *testing.T) {
		code := `
	let xhr1 =new XMLHttpRequest()

	xhr1.open('GET', 'http://example.com')
	xhr1.send()
    xhr1.send("123")
    xhr1.addEventListener('load', function () {
      console.log(this.response)
    })
   `
		prog, err := ssaapi.Parse(code, ssaapi.WithLanguage(ssaapi.JS))
		if err != nil {
			t.Fatal("prog parse error", err)
		}
		prog.Show()
		results, err := prog.SyntaxFlowWithError("XMLHttpRequest().open")
		for _, result := range results {
			// 获取所有call被调用的地方
			for _, called := range result.GetCalledBy() {
				//获取参数
				method := called.GetCallArgs()[0]
				fmt.Printf("method: %v\n", method.String())
				require.Equal(t, method.String(), "\"GET\"")
				path := called.GetCallArgs()[1]
				fmt.Printf("path: %v\n", path.String())
				require.Equal(t, path.String(), "\"http://example.com\"")
			}
		}
	})

	t.Run("simple post request", func(t *testing.T) {
		code := `
	const data1 = {
       name: 'job',
       age: '12',
    }
    let xhr1 = new XMLHttpRequest()
    xhr1.open('POST', 'http://example1.com')
    const usp = new URLSearchParams(data)
    const query = usp.toString()
    xhr1.setRequestHeader('Content-type', 'application/x-www-form-urlencoded')
    xhr1.send(query)
    xhr1.addEventListener('load', function () {
        console.log(this.response)
    })

const data2 = {
       name: 'job',
       age: '12',
    }
    let xhr2 = new XMLHttpRequest()
    xhr2.open('POST', 'http://example2.com')
    const usp = new URLSearchParams(data)
    const query = usp.toString()
    xhr2.setRequestHeader('Content-type', 'application/x-www-form-urlencoded')
    xhr2.send(query)
    xhr2.addEventListener('load', function () {
        console.log(this.response)
    })

   `
		prog, err := ssaapi.Parse(code, ssaapi.WithLanguage(ssaapi.JS))
		if err != nil {
			t.Fatal("prog parse error", err)
		}
		prog.Show()
		xmlRequest, err := prog.SyntaxFlowWithError("XMLHttpRequest()")
		open := xmlRequest.Ref("open")
		if len(open) != 0 {
			for _, result := range open {
				// 获取所有call被调用的地方
				for _, called := range result.GetCalledBy() {
					//获取参数
					fmt.Println("=====================================================")
					method := called.GetCallArgs()[0]
					fmt.Printf("method: %v\n", method.String())
					require.Equal(t, method.String(), "\"POST\"")
					path := called.GetCallArgs()[1]
					fmt.Printf("path: %v\n", path.String())
					require.Contains(t, path.String(), "http://example")
				}
				// 获取post的data
			}
		}

	})

}

func TestJs_Ajax(t *testing.T) {
	code := `$.ajax({ //统计访问量
    url:'/foot_action!getCount.action',
    type: 'POST',
    dataType: 'json',
    cache:false,
    data: {url:window.location.href},
    timeout: 5000,
    error: function(){
    },
    success: function(result){
     $("#fwls").html(result.count);
        $("#fwl").html(result.count1);
    }
 });`
	prog, err := ssaapi.Parse(code, ssaapi.WithLanguage(ssaapi.JS))
	if err != nil {
		panic(err)
	}
	withError, err := prog.SyntaxFlowWithError("$.ajax()")
	if err != nil {
		panic(err)
	}
	for _, value := range withError {
		operator, err := value.GetCallActualParams()
		if err != nil {
			panic(err)
		}
		members, err := operator.GetMembers()
		if err != nil {
			panic(err)
		}
		match, valueOperator, err := members.ExactMatch("url")
		if err != nil {
			panic(err)
		}
		assert.Equal(t, match, true, "match is false, except true")
		assert.Equal(t, len(valueOperator.GetNames()), 4, "valueOperator number not match")
	}
}
func TestAjax_post(t *testing.T) {
	code := `$.post({
  url: 'https://jsonplaceholder.typicode.com/posts',
  contentType: 'application/json',
  data: JSON.stringify(formData),
  success: function(response) {
    // ...
  },
  error: function(xhr, status, error) {
    // ...
  }
});
`
	prog, err := ssaapi.Parse(code, ssaapi.WithLanguage(ssaapi.JS))
	if err != nil {
		panic(err)
	}
	withError, err := prog.SyntaxFlowWithError("$.post()")
	if err != nil {
		panic(err)
	}
	params, err := withError.GetCallActualParams()
	if err != nil {
		panic(err)
	}
	members, err := params.GetMembers()
	if err != nil {
		panic(err)
	}
	match, operator, err := members.ExactMatch("url")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, match, true, "match is false, except true")
	assert.Equal(t, len(operator.GetNames()), 4, "valueOperator number not match")
}

func Test_JS_Fetch(t *testing.T) {
	code := `fetch('url')
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json(); 
  })
  .then(data => {
    console.log(data);
  })
  .catch(error => {
    console.error('There has been a problem with your fetch operation:', error);
  });
`
	prog, err := ssaapi.Parse(code, ssaapi.WithLanguage(ssaapi.JS))
	if err != nil {
		panic(err)
	}
	withError, err := prog.SyntaxFlowWithError("fetch()")
	if err != nil {
		panic(err)
	}
	params, err := withError.GetCallActualParams()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, len(params.GetNames()), 2, fmt.Sprintf("not match,except 1,match %v", len(params.GetNames())))
}

func Test_JS_JQuery(t *testing.T) {}

func Test_JS_Axios(t *testing.T) {}