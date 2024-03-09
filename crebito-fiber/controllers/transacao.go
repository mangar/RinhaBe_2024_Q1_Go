package controllers

import (
	"crebito-fiber/helpers"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TransacaoInput struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type TransacaoResult struct {
	Limite       int    `json:"limite"`
	Saldo        int 	`json:"saldo"`
	StatusCode int `json:"-"`
	StatusMessage string `json:"-"`	
}

type TransacaoController struct {
	idCliente 	int
	input		TransacaoInput
	c 			*fiber.Ctx
}

func NewTransacaoController(c *fiber.Ctx) (*TransacaoController) {
	return &TransacaoController{c:c}
}


func (e *TransacaoController) Run() TransacaoResult {

	idParam := e.c.Params("ID", "0")

	idCliente, err := strconv.Atoi(idParam)
	if err != nil  {
		logrus.Info("[TransacaoController.Run] idParam !int:", idParam)
		return TransacaoResult{StatusCode: 422}
	}

	if idCliente < 1 || idCliente > 5 {
		return TransacaoResult{StatusCode: 404}
	}
	e.idCliente = idCliente


	var input TransacaoInput

	err = e.c.BodyParser(&input)
	if err != nil {
		logrus.Info("[TransacaoController.Run] BodyParser")
		return TransacaoResult{StatusCode: 422}
	}

	if input.Value < 1 {
		logrus.Info("[TransacaoController.Run] Valor Invalido")
		return TransacaoResult{StatusCode: 422, StatusMessage: "[TransacaoController.Run] Valor Invalido"}
	}

	if len(input.Description) < 1 || len(input.Description) > 10 {
		logrus.Info("[TransacaoController.Run] Descricao deve ter entre 1 e 10 caracteres")
		return TransacaoResult{StatusCode: 422, StatusMessage: "[TransacaoController.Run] Descricao deve ter entre 1 e 10 caracteres"}
	}

	if input.Type != "c" && input.Type != "d" {
		logrus.Info("[TransacaoController.Run] Tipo de transacao invalida: " + input.Type)
		return TransacaoResult{StatusCode: 422, StatusMessage: "[TransacaoController.Run] Tipo de transacao invalida: " + input.Type}
	}

	e.input = input

	ssaldo, llimite, err := e.runProc()
	if err != nil {
		return TransacaoResult{StatusCode: 422, StatusMessage: err.Error()}
	}

	return TransacaoResult{StatusCode: 200, Limite: llimite, Saldo: ssaldo}

}


func (e *TransacaoController) runProc() (int, int, error) {

	var ssaldo int
	var llimite int
	var err error
	if e.input.Type == "c" {
		err = helpers.GetDBConnection().QueryRow(e.c.Context(), 
        "SELECT * FROM func_credito($1,$2,$3,$4,$5);", 
		uuid.New().String(), e.idCliente, e.input.Value, e.input.Description, time.Now()).Scan(&ssaldo, &llimite)

	} else {
		err = helpers.GetDBConnection().QueryRow(e.c.Context(), 
        "SELECT * FROM func_debito($1,$2,$3,$4,$5);", 
		uuid.New().String(), e.idCliente, e.input.Value, e.input.Description, time.Now()).Scan(&ssaldo, &llimite)
	}
            
    if err != nil {
		logrus.Debug("[TransacaoController.run] Erro inesperado ...", err.Error())
		return ssaldo, llimite, err
	}
	return ssaldo, llimite, nil
}


