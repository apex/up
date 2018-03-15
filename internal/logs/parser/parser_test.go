package parser

import (
	"testing"
)

// TODO: precedence...
// TODO: byte size literals
// TODO: support literals: `method = GET`, `ip = 70.*` etc
// TODO: error tests
// TODO: test error messages
// TODO: add "starts with" / "ends with" to compliment "contains"?
// TODO: document best practices for logging json from an app

var cases = []struct {
	Input  string
	Output string
}{
	{`production`, `{ $.fields.stage = "production" }`},
	{`development`, `{ $.fields.stage = "development" }`},
	{`staging`, `{ $.fields.stage = "staging" }`},
	{`method = "GET"`, `{ $.fields.method = "GET" }`},
	{`debug`, `{ ($.level = "debug") }`},
	{`info`, `{ ($.level = "info") }`},
	{`warn`, `{ ($.level = "warn") }`},
	{`error`, `{ ($.level = "error") }`},
	{`fatal`, `{ ($.level = "fatal") }`},
	{`not info`, `{ !(($.level = "info")) }`},
	{`not error or fatal`, `{ !(($.level = "error") || ($.level = "fatal")) }`},
	{`!info`, `{ !($.level = "info") }`},
	{`level = "info"`, `{ $.level = "info" }`},
	{`message = "user signin"`, `{ $.message = "user signin" }`},
	{`email = "tj@apex.sh"`, `{ $.fields.email = "tj@apex.sh" }`},
	{`status = 0`, `{ $.fields.status = 0 }`},
	{`status = 0.123`, `{ $.fields.status = 0.123 }`},
	{`status = .123`, `{ $.fields.status = 0.123 }`},
	{`status = 200`, `{ $.fields.status = 200 }`},
	{`price = 1.95`, `{ $.fields.price = 1.95 }`},
	{`price == 1.95`, `{ $.fields.price = 1.95 }`},
	{`price > 1.95`, `{ $.fields.price > 1.95 }`},
	{`price < 1.95`, `{ $.fields.price < 1.95 }`},
	{`price >= 1.95`, `{ $.fields.price >= 1.95 }`},
	{`price <= 1.95`, `{ $.fields.price <= 1.95 }`},
	{`price != 1.95`, `{ $.fields.price != 1.95 }`},
	{`!enabled`, `{ !$.fields.enabled }`},
	{`!   enabled`, `{ !$.fields.enabled }`},
	{`foo = 1 || bar = 2`, `{ $.fields.foo = 1 || $.fields.bar = 2 }`},
	{`foo = 1 && bar = 2`, `{ $.fields.foo = 1 && $.fields.bar = 2 }`},
	{`foo = 1 or bar = 2`, `{ $.fields.foo = 1 || $.fields.bar = 2 }`},
	{`foo = 1 and bar = 2`, `{ $.fields.foo = 1 && $.fields.bar = 2 }`},
	{`foo = 1 bar = 2`, `{ $.fields.foo = 1 && $.fields.bar = 2 }`},
	{`foo.bar.baz = 1`, `{ $.fields.foo.bar.baz = 1 }`},
	{`level = "error" and (duration >= 500 or duration = 0)`, `{ $.level = "error" && ($.fields.duration >= 500 || $.fields.duration = 0) }`},
	{`level = "error" (duration >= 500 or duration = 0)`, `{ $.level = "error" && ($.fields.duration >= 500 || $.fields.duration = 0) }`},
	{`cart.total = 15.99`, `{ $.fields.cart.total = 15.99 }`},
	{`user.name contains "obi"`, `{ $.fields.user.name = "*obi*" }`},
	{`user in ("Tobi")`, `{ ($.fields.user = "Tobi") }`},
	{`pet.age in (1, 2, 3)`, `{ ($.fields.pet.age = 1 || $.fields.pet.age = 2 || $.fields.pet.age = 3) }`},
	{`user in ("Tobi", "Loki", "Jane")`, `{ ($.fields.user = "Tobi" || $.fields.user = "Loki" || $.fields.user = "Jane") }`},
	{`user.name in ("Tobi", "Loki", "Jane")`, `{ ($.fields.user.name = "Tobi" || $.fields.user.name = "Loki" || $.fields.user.name = "Jane") }`},
	{`not user.admin`, `{ !($.fields.user.admin) }`},
	{`not user.role in ("Admin", "Moderator")`, `{ !(($.fields.user.role = "Admin" || $.fields.user.role = "Moderator")) }`},
	{`user.role not in ("Admin", "Moderator")`, `{ !(($.fields.user.role = "Admin" || $.fields.user.role = "Moderator")) }`},
	{`not level = "error" or level = "fatal"`, `{ !($.level = "error" || $.level = "fatal") }`},
	{`cart.products[0] = "something"`, `{ $.fields.cart.products[0] = "something" }`},
	{`cart.products[0].price = 15.99`, `{ $.fields.cart.products[0].price = 15.99 }`},
	{`cart.products[0][1].price = 15.99`, `{ $.fields.cart.products[0][1].price = 15.99 }`},
	{`cart.products[0].items[1].price = 15.99`, `{ $.fields.cart.products[0].items[1].price = 15.99 }`},
	{`user.name in ("Tobi", "Loki") and status >= 500`, `{ ($.fields.user.name = "Tobi" || $.fields.user.name = "Loki") && $.fields.status >= 500 }`},
	{`method in ("POST", "PUT") and ip = "207.*" and status = 200 and duration >= 50`, `{ ($.fields.method = "POST" || $.fields.method = "PUT") && $.fields.ip = "207.*" && $.fields.status = 200 && $.fields.duration >= 50 }`},
	{`method in ("POST", "PUT") ip = "207.*" status = 200 duration >= 50`, `{ ($.fields.method = "POST" || $.fields.method = "PUT") && $.fields.ip = "207.*" && $.fields.status = 200 && $.fields.duration >= 50 }`},
	{`size > 1kb`, `{ $.fields.size > 1024 }`},
	{`size > 2kb`, `{ $.fields.size > 2048 }`},
	{`size > 1.5mb`, `{ $.fields.size > 1572864 }`},
	{`size > 100b`, `{ $.fields.size > 100 }`},
	{`duration > 100ms`, `{ $.fields.duration > 100 }`},
	{`duration > 1s`, `{ $.fields.duration > 1000 }`},
	{`duration > 4.5s`, `{ $.fields.duration > 4500 }`},
	{`"User Login"`, `{ $.message = "User Login" }`},
	{`"User*"`, `{ $.message = "User*" }`},
	{`"Signup" or "Signin"`, `{ $.message = "Signup" || $.message = "Signin" }`},
	{`"User Login" method = "GET"`, `{ $.message = "User Login" && $.fields.method = "GET" }`},
	{`method = GET`, `{ $.fields.method = "GET" }`},
	{`method in (GET, HEAD, OPTIONS)`, `{ ($.fields.method = "GET" || $.fields.method = "HEAD" || $.fields.method = "OPTIONS") }`},
	{`name = tj`, `{ $.fields.name = "tj" }`},
	{`method = GET path = /account/billing`, `{ $.fields.method = "GET" && $.fields.path = "/account/billing" }`},
	{`cart.products[0].name = ps4`, `{ $.fields.cart.products[0].name = "ps4" }`},
}

func TestParse(t *testing.T) {
	for _, c := range cases {
		t.Logf("parsing %q", c.Input)
		n, err := Parse(c.Input)

		if err != nil {
			t.Errorf("error parsing %q: %s", c.Input, err)
			continue
		}

		if n.String() != c.Output {
			t.Errorf("\n\ntext: %s\nwant: %s\n got: %s\n\n", c.Input, c.Output, n.String())
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse(`user.name in ("Tobi", "Loki", "Jane")`)
	}
}
