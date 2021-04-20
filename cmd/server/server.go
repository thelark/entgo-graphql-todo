// Copyright 2021 CloudWay
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"net/http"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"entgqlgo/ent"
	"entgqlgo/ent/migrate"
	"entgqlgo/graph"

	_ "entgqlgo/ent/runtime"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open(
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_fk=1",
	)
	if err != nil {
		log.Fatal("opening ent client", err)
	}
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Fatal("running schema migration", err)
	}
	srv := handler.NewDefaultServer(graph.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	http.Handle("/",
		playground.Handler("Todo", "/query"),
	)
	http.Handle("/query", srv)

	log.Println("listening on", ":9898")
	if err := http.ListenAndServe(":9898", nil); err != nil {
		log.Println("http server terminated", err)
	}
}
