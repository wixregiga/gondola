// Package bootstrap3 implements some helper functions intended
// to be used with the Bootstrap front-end framework, version 3.
// See http://twitter.github.io/bootstrap/ for more details.
//
// This package defines the "bootstrap" asset, which serves bootstrap
// from http://www.bootstrapcdn.com. It receives a single argument with the
// desired bootstrap version. Only version 3 is supported. e.g.
//
//  bootstrap: 3.0.0
//
// This asset also supports the following options:
//
//  nojs (bool): disables loading bootstrap's javascript library
//  e.g. bootstrap|nojs: 3.0.0
//
// See gondola/template and gondola/template/assets for more information
// about template functions and the assets pipeline.
//
// Importing this package will also register FormRenderer as the default
// gondola/form renderer and PaginatorRenderer as the default
// gondola/html/paginator renderer.
package bootstrap3
