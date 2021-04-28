// Package articles implements an app for displaying articles from files.
//
// This application is intended to display a collection of articles where each
// article is represented by a file. The file must have one of the following
// extensions:
//
//  .html, .txt, .md
//
// To load articles written in markdown (.md extension), your application must
// import gondola/template/markdown, like e.g.
//
//  import (
//	    _ "gondola/template/markdown"
//  )
//
// Each file is composed of first the article text and then a set of properties
// separated by one empty line. Each property takes the form:
//
//  [name] = value
//
// Where value must a valid value for the property type. The currently parsed
// properties are:
//
//  - title (required): Indicates the article title and can be used to generate the
//	article URL when no slug is specified.
//  - id (optional): The article id, used to reverse the article. If it's not present,
//	the article filename without extension is used as its id.
//  - slug (optional): The article slug. If present, it overrides the generated value
//	from the title.
//  - synopsis (optional): A small text about the article shown in article listings.
//  - updated (optional): Indicates an update to the file.
//  - priority (optional): When listing the articles, the ones with lower priority
//	are shown first.
//
// The title, slug and updated field might appear multiple times. It's recommended that if you
// change the title or the slug of an article, you do so by adding a new property BEFORE the
// previous one without deleting the old one. This allows the articles app to redirect users
// from the old URL to the new one.
//
// This package also adds a template function named reverse_article. It can be used to find the
// URL of an article from its id. e.g.
//
//      {{ reverse_article "article-id" }}
//
// The typical usage of this application is as follows:
//
//      articlesApp := articles.New()
//      if _, err := articlesApp.LoadDir(pathutil.Relative("articles")); err != nil {
//          panic(err)
//      }
//
//      myapp.Include("/articles/", articlesApp.App, "articles-base.html")
//
package articles
