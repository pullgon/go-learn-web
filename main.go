package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

var router = mux.NewRouter()

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func main() {

	router.HandleFunc("/", homeHandler).Methods("GET", "POST").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("get").Name("about")

	// article
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articleCreateHandler).Methods("GET").Name("articles.create")

	// 404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件 text/html
	router.Use(forceHTMLMiddleware)
	// router优化：去掉结尾的"/"
	// router.Use(removeTrailingSlash)

	// 获取url
	homeURL := router.Get("home")
	fmt.Println("homeURL：", homeURL.GetName())
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL：", articleURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}

func articleCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<title>创建文章 —— 我的技术博客</title>
		</head>
		<body>
			<form action="%s?test=data" method="post">
				<p><input type="text" name="title"></p>
				<p><textarea name="body" cols="30" rows="10"></textarea></p>
				<p><button type="submit">提交</button></p>
			</form>
		</body>
		</html>
	`
	url, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, url)
}

func removeTrailingSlash(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		handler.ServeHTTP(w, r)
	})
}

func forceHTMLMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置标头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		// 继续处理请求
		handler.ServeHTTP(w, r)
	})
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "文章ID：%s", id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	// 标题不能为空，且要大于两个字符，且小于 40 个字符
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if len(title) < 3 || len(title) > 40 {
		errors["title"] = "标题长度需介于3-40"
	}

	// 内容不能为空，且要大于 10 个字符
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if len(body) <= 10 {
		errors["body"] = "内容长度不能少于10个字符"
	}

	if len(errors) == 0 {
		fmt.Fprintf(w, "验证通过！<br>")
		fmt.Fprintf(w, "title: %s<br>", title)
		fmt.Fprintf(w, "title len: %d<br>", len(title))
		fmt.Fprintf(w, "body: %s<br>", body)
		fmt.Fprintf(w, "body len: %d<br>", len(body))
	} else {
		errHtml := `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>创建文章 —— 我的技术博客</title>
    <style type="text/css">.error {color: red;}</style>
</head>
<body>
    <form action="{{ .URL }}" method="post">
        <p><input type="text" name="title" value="{{ .Title }}"></p>
        {{ with .Errors.title }}
        <p class="error">{{ . }}</p>
        {{ end }}
        <p><textarea name="body" cols="30" rows="10">{{ .Body }}</textarea></p>
        {{ with .Errors.body }}
        <p class="error">{{ . }}</p>
        {{ end }}
        <p><button type="submit">提交</button></p>
    </form>
</body>
</html>
`
		//fmt.Fprintf(w, "有错误发生, error值为：%v", errors)
		formUrl, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{Title: title, Body: body, URL: formUrl, Errors: errors}
		t, err := template.New("create-form").Parse(errHtml)
		if err != nil {
			panic(err)
		}
		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>hello, 欢迎来到 goBlog</h1>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系"+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}
