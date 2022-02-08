package query

import (
	"github.com/nbs-go/nsql/schema"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
	"time"
)

func TestSchemaBuilder(t *testing.T) {
	type Customer struct {
		CreatedAt time.Time `db:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt"`
		Id        int64     `db:"id"`
		FullName  string    `db:"fullName"`
		Gender    string    `db:"gender"`
	}

	s := schema.New(schema.FromModelRef(Customer{}))
	sb := Schema(s)

	// Test #1
	test_utils.CompareString(t, "FIND BY PRIMARY KEY", sb.FindByPK(),
		`SELECT "Customer"."createdAt", "Customer"."updatedAt", "Customer"."id", "Customer"."fullName", "Customer"."gender" FROM "Customer" WHERE "Customer"."id" = ?`)

	// Test #2
	test_utils.CompareString(t, "INSERT", sb.Insert(),
		`INSERT INTO "Customer"("createdAt", "updatedAt", "fullName", "gender") VALUES (:createdAt, :updatedAt, :fullName, :gender)`)

	// Test #3
	test_utils.CompareString(t, "UPDATE", sb.Update(),
		`UPDATE "Customer" SET "createdAt" = :createdAt, "updatedAt" = :updatedAt, "fullName" = :fullName, "gender" = :gender WHERE "id" = :id`)

	// Test #4
	test_utils.CompareString(t, "DELETE", sb.Delete(),
		`DELETE FROM "Customer" WHERE "id" = ?`)

	// Test #5
	test_utils.CompareString(t, "COUNT", sb.Count(Like(Column("fullName"))),
		`SELECT COUNT("Customer"."id") AS "count" FROM "Customer" WHERE "Customer"."fullName" LIKE ?`)
}
