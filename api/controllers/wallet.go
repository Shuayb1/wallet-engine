package controllers

import (
	"opay/api/auth"
	models "opay/api/models"
	responses "opay/api/responses"
	utils "opay/api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//CreateWallet method to create wallet with body parameter nickname,email,password
func (server *Server) CreateWallet(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	wallet := models.Wallet{}
	err = json.Unmarshal(body, &wallet)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	wallet.Prepare()
	err = wallet.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	walletCreated, err := wallet.SaveWallet(server.DB)

	if err != nil {

		formattedError := utils.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, walletCreated.WalletID))
	responses.JSON(w, http.StatusCreated, walletCreated)
}

//GetWallets method to get all wallet
func (server *Server) GetWallets(w http.ResponseWriter, r *http.Request) {

	wallet := models.Wallet{}

	wallets, err := wallet.FindAllWallets(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, wallets)
}

//GetWallet method to get wallet by id in case with parameter route
func (server *Server) GetWallet(w http.ResponseWriter, r *http.Request) {
	//CASE USING FORM DATA ID
	// var id = r.FormValue("id")
	// uid, err := strconv.ParseUint(id, 10, 32)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	// wallet := models.Wallet{}
	// walletGotten, err := wallet.FindWalletByID(server.DB, uint32(uid))
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	// responses.JSON(w, http.StatusOK, walletGotten)

	//==========================
	//CASE USING ROUTES USERS/{ID}
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	wallet := models.Wallet{}
	walletGotten, err := wallet.FindWalletByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, walletGotten)
}

//UpdateWallet method to update wallet with route parameter id
func (server *Server) CreditWallet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	wallet := models.Wallet{}
	var dat map[string]interface{}


	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &dat)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	amount :=  dat["CreditAmount"].(float64)
	
	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if tokenID != uint32(uid) {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }


	updatedWallet, err := wallet.CreditAWallet(server.DB, uint32(uid), uint32(amount))

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}


	responses.JSON(w, http.StatusOK, updatedWallet)
}

func (server *Server) DebitWallet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var dat map[string]interface{}
	wallet := models.Wallet{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if tokenID != uint32(uid) {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	amount :=  dat["DeditAmount"].(float64)
	
	updatedWallet, err := wallet.DebitAWallet(server.DB, uint32(uid), uint32(amount))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedWallet)
}

func (server *Server) ActivateAWallet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	wallet := models.Wallet{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if tokenID != 0 && tokenID != uint32(uid) {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	updatedWallet, err := wallet.ActivateAWallet(server.DB, uint32(uid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedWallet)
}

func (server *Server) DeactivateAWallet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	wallet := models.Wallet{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if tokenID != 0 && tokenID != uint32(uid) {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	updatedWallet, err := wallet.DeactivateAWallet(server.DB, uint32(uid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedWallet)
}

//DeleteWallet method to delete wallet from db
func (server *Server) DeleteWallet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	wallet := models.Wallet{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = wallet.DeleteAWallet(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
