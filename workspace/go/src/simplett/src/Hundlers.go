package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

func GetUsers(w http.ResponseWriter,r *http.Request){
	db := GetDB()
	defer db.Close()
	basicusers := PgGetUsers(db)

	//1、在头部输出内容格式json格式告知客户端
	//2、设置http状态码
	w.Header().Set("Content-Type","application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//结构体转换成json格式并输出
	if err := json.NewEncoder(w).Encode(basicusers); err != nil {
		panic(err)
	}

}

func UserCreate(w http.ResponseWriter,r *http.Request){
	db := GetDB()
	defer db.Close()
	var basicuser BasicUser
	body,err := ioutil.ReadAll(io.LimitReader(r.Body,1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &basicuser); err != nil {
		//设置格式和状态吗
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	basicuser.Type = "user"
	PgUserCreate(db,basicuser)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	basicuser.Id = PgGetUserIdByName(db,basicuser.Name)
	if err := json.NewEncoder(w).Encode(basicuser); err != nil {
		panic(err)
	}
}

func GetRelationships(w http.ResponseWriter,r *http.Request) {
	vars := mux.Vars(r)
	id,err := strconv.Atoi(vars["UserId"])
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	fmt.Println(reflect.TypeOf(id).String())
	db := GetDB()
	liked  := GetLiked(db,id)
	fmt.Println(liked)
	relations := PgGetRelationships(db,id)
	if err := json.NewEncoder(w).Encode(relations); err != nil {
		panic(err)
	}
}

func PutRelationships(w http.ResponseWriter,r *http.Request)  {

	//获取参数id 和otherId并转换为int类型
	vars := mux.Vars(r)
	id,err := strconv.Atoi(vars["UserId"])
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	otherId,err := strconv.Atoi(vars["OtherUserId"])
	if err != nil {
		panic(err)
	}
	fmt.Println(otherId)

	//获取数据库连接
	db := GetDB()
	defer db.Close()

	//获取json数据并转换为Relation类型
	var relation Relation
	body,err := ioutil.ReadAll(io.LimitReader(r.Body,1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &relation); err != nil {
		//设置格式和状态码
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//根据读到的数据更新relation，并持久化到pg
	relation.Type = "relationship"
	relation.Id = otherId
	relation = PgPutRelationships(db,id,otherId,relation)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(relation); err != nil {
		panic(err)
	}

}