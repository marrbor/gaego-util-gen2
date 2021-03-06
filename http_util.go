/*
 * Copyright (c) 2019 Masahide Matsumoto
 * -*- coding:utf-8 -*-
 *
 * Web API を構築するためのユーティリティ
 *
 */
package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Request

// RequestToParams は、request にあるボディー（JSON）を params でポインタ渡しされた go struct に変換します。
func RequestToParams(r *http.Request, params interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(&params)
}

// Response

// OK は、200 OK のヘッダだけを返します。
func OK(w http.ResponseWriter) {
	w.WriteHeader(200)
}

// TextResponse は、200 OKで指定テキストを返します。
func TextResponse(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(text)); err != nil {
		InternalServerError(w, err)
	}
}

// JSONResponse は、200 OKで JSON オブジェクトを返します。エンコードに失敗したら 500 エラーを返します。
func JSONResponse(w http.ResponseWriter, data interface{}) {
	if data == nil {
		OK(w)
		return
	}
	j, err := json.Marshal(data)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		InternalServerError(w, err)
	}
}

// エラーを返す
func errResponse(w http.ResponseWriter, code int, err error) {
	msg := http.StatusText(code)
	if err != nil {
		msg = msg + " " + err.Error()
	}
	http.Error(w, msg, code)
}

// BadRequest は、400エラーを返します。
func BadRequest(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusBadRequest, err)
}

// Unauthorized は、401エラーを返します。
func Unauthorized(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusUnauthorized, err)
}

// Forbidden は、403エラーを返します。
func Forbidden(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusForbidden, err)
}

// NotFound は、404エラーを返します。
func NotFound(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusNotFound, err)
}

// MethodNotAllowed は、405エラーを返します。
func MethodNotAllowed(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusMethodNotAllowed, err)
}

// InternalServerError は、500エラーを返します。
func InternalServerError(w http.ResponseWriter, err error) {
	errResponse(w, http.StatusInternalServerError, err)
}

// StartServer は環境変数「PORT」の指定または指定されたポート番号で待受を開始し、ポート番号を返します。
func StartServer(port int32, handler http.Handler) error {
	portStr := GetPort()
	if portStr == "" {
		portStr = fmt.Sprintf("%d", port)
	}

	// ポート番号チェック
	if _, err := strconv.ParseInt(portStr, 10, 32); err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", portStr), handler)
}
