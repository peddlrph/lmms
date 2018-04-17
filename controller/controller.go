// Package controller loads the routes for each of the controllers.
package controller

import (
	"github.com/blue-jay/blueprint/controller/about"
	"github.com/blue-jay/blueprint/controller/cash"
	"github.com/blue-jay/blueprint/controller/code"
	"github.com/blue-jay/blueprint/controller/debug"
	"github.com/blue-jay/blueprint/controller/expenses"
	"github.com/blue-jay/blueprint/controller/home"
	"github.com/blue-jay/blueprint/controller/load"
	"github.com/blue-jay/blueprint/controller/login"
	"github.com/blue-jay/blueprint/controller/message"
	"github.com/blue-jay/blueprint/controller/mobile_ip"
	"github.com/blue-jay/blueprint/controller/notepad"
	"github.com/blue-jay/blueprint/controller/register"
	"github.com/blue-jay/blueprint/controller/smartmoney"
	"github.com/blue-jay/blueprint/controller/smartpadala"
	"github.com/blue-jay/blueprint/controller/static"
	"github.com/blue-jay/blueprint/controller/status"
	"github.com/blue-jay/blueprint/controller/summary"
	"github.com/blue-jay/blueprint/controller/transaction"
)

// LoadRoutes loads the routes for each of the controllers.
func LoadRoutes() {
	about.Load()
	debug.Load()
	register.Load()
	login.Load()
	home.Load()
	static.Load()
	status.Load()
	notepad.Load()
	cash.Load()
	load.Load()
	smartmoney.Load()
	smartpadala.Load()
	transaction.Load()
	code.Load()
	summary.Load()
	expenses.Load()
	message.Load()
	mobile_ip.Load()
}
