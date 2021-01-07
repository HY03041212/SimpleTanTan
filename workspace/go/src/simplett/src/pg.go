package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "simplett"
)

func GetDB() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func PgGetUsers(db *sql.DB) BasicUsers {
	var basicusers BasicUsers
	var basicuser BasicUser
	sqlStatement := `SELECT id,name,type FROM users`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&basicuser.Id, &basicuser.Name, &basicuser.Type)
		basicusers = append(basicusers, basicuser)
	}
	return basicusers
}

func PgGetUserIdByName(db *sql.DB, name string) int {
	var id int
	sqlStatement := "SELECT id FROM users where name = $1"
	row := db.QueryRow(sqlStatement, name)
	row.Scan(&id)
	return id
}

func PgUserCreate(db *sql.DB, basicUser BasicUser) {
	stmt, err := db.Prepare("insert into users(name,type) values($1,$2)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(basicUser.Name, basicUser.Type)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("insert success")
	}
}

func PgGetRelationships(db *sql.DB, id int) Relations {
	var relations Relations
	var relation Relation
	var liked []sql.NullInt64
	var matched []sql.NullInt64
	var disliked []sql.NullInt64
	sqlStatement := "select liked,matched,disliked from users where id = $1"
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(pq.Array(&liked), pq.Array(&matched), pq.Array(&disliked))
	if err != nil {
		fmt.Println(err.Error())
		return relations
	}
	fmt.Println(liked)
	fmt.Println(matched)
	fmt.Println(disliked)

	//遍历数组，把数据封装到relations结构体切片类型里面
	for _, userid := range liked {
		relation.Id = int(userid.Int64)
		relation.State = "liked"
		relation.Type = "relationship"
		relations = append(relations, relation)
	}
	for _, userid := range matched {
		relation.Id = int(userid.Int64)
		relation.State = "matched"
		relation.Type = "relationship"
		relations = append(relations, relation)
	}
	for _, userid := range disliked {
		relation.Id = int(userid.Int64)
		relation.State = "disliked"
		relation.Type = "relationship"
		relations = append(relations, relation)
	}
	fmt.Println(relations)
	return relations
}

func GetLiked(db *sql.DB, id int) (liked []sql.NullInt64) {
	sqlStatement := "select liked from users where id = $1"
	row := db.QueryRow(sqlStatement, id)
	row.Scan(pq.Array(&liked))
	return liked
}

func GetMatched(db *sql.DB, id int) (matched []sql.NullInt64) {
	sqlStatement := "select liked from users where id = $1"
	row := db.QueryRow(sqlStatement, id)
	row.Scan(pq.Array(&matched))
	return matched
}

func GetDisliked(db *sql.DB, id int) (disliked []sql.NullInt64) {
	sqlStatement := "select liked from users where id = $1"
	row := db.QueryRow(sqlStatement, id)
	row.Scan(pq.Array(&disliked))
	return disliked
}

func PgUpdateMatched(db *sql.DB, id int, otherid int) {
	var matched []sql.NullInt64
	var otherMatched []sql.NullInt64
	matched = GetMatched(db, id)
	var  tmpid sql.NullInt64
	tmpid.Int64 = int64(otherid)
	matched = append(matched, tmpid)
	var matched2 []int64
	for _,onematched := range matched {
		matched2 = append(matched2,onematched.Int64)
	}

	otherMatched = GetMatched(db, otherid)
	var  othertmpid sql.NullInt64
	othertmpid.Int64 = int64(id)
	otherMatched = append(otherMatched, othertmpid)
	var othermatched2 []int64
	for _,onematched := range otherMatched {
		matched2 = append(othermatched2,onematched.Int64)
	}

	sqlStatement := "update users set matched = $1 where id = $2"
	_, err := db.Exec(sqlStatement, pq.Array(matched2), id)
	if err != nil {
		panic(err)
	}
	_, err2 := db.Exec(sqlStatement, pq.Array(othermatched2), otherid)
	if err2 != nil {
		panic(err2)
	}
}

func PgPutRelationships(db *sql.DB, id int, otherid int, relation Relation) Relation {
	var feel []sql.NullInt64
	var otherLiked []sql.NullInt64
	if relation.State == "liked" {
		feel = GetLiked(db, id)
		var  tmpid sql.NullInt64
		tmpid.Int64 = int64(otherid)
		feel = append(feel,tmpid)
		var feel2 []int64
		for _,onefeel := range feel {
			feel2 = append(feel2,onefeel.Int64)
		}
		sqlStatement := "update users set liked = $1 where id = $2"
		_, err := db.Exec(sqlStatement, pq.Array(feel2),id)
		if err != nil {
			panic(err)
		}

		//A，B两人处理match问题：遍历B的liked数组，查询是否有A.id
		otherLiked = GetLiked(db, otherid)
		for _, likedId := range otherLiked {
			if likedId.Int64 == int64(id) {
				PgUpdateMatched(db, id, otherid)
				relation.State = "matched"
				break
			}
		}

	} else {
		feel = GetDisliked(db, id)
		var  tmpid sql.NullInt64
		tmpid.Int64 = int64(otherid)
		feel = append(feel,tmpid)
		fmt.Println(feel)
		var feel2 []int64
		for _,onefeel := range feel {
			feel2 = append(feel2,onefeel.Int64)
		}
		fmt.Println(feel2)
		sqlStatement := "update users set disliked = $1 where id = $2"
		_, err := db.Exec(sqlStatement, pq.Array(feel2),id)
		if err != nil {
			panic(err)
		}
	}
	return relation
}
