import (
	"encoding/json"
	"net/http"
)

//Coin balance params
 type CoinBalanceParams struct {
	Usernname string
 }

 //Coin base balance response
 type CoinBalanceResponse struct {
	Usernname string
	Balance int64
 }

 //Error responce

 type ErrorResponse struct {
	Code int
	Message string
 }