package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var router = mux.NewRouter()

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
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "解析错误，请提供正确的数据！")
		return
	}
	title := r.PostForm.Get("title")

	// PostForm 存储了 post、put 参数，在使用之前需要调用 ParseForm 方法
	fmt.Fprintf(w, "PostForm：%v<br>", r.PostForm)
	// 存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。
	fmt.Fprintf(w, "form：%v<br>", r.Form)
	fmt.Fprintf(w, "title：%v<br>", title)

	// FormValue 和 PostFormValue无需使用.ParseForm()
	fmt.Fprintf(w, "r.Form 中 title 的值为: %v <br>", r.FormValue("title"))
	fmt.Fprintf(w, "r.PostForm 中 title 的值为: %v <br>", r.PostFormValue("title"))
	fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
	fmt.Fprintf(w, "r.PostForm 中 test 的值为: %v <br>", r.PostFormValue("test"))

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
