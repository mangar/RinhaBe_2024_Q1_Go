package controllers

import (
	"crebito-fiber/helpers"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ExtratoResult struct {
	Saldo ExtratoSaldoResult `json:"saldo"`
	UltimasTransacoes []ExtratoTransacaoResult `json:"ultimas_transacoes"`
	StatusCode int `json:"-"`
	StatusMessage string `json:"-"`
}
type ExtratoSaldoResult struct {
	Total  int        		`json:"total"`
	DataExtrato time.Time 	`json:"data_extrato"`
	Limite int        		`json:"limite"`	
}
type ExtratoTransacaoResult struct {
	Valor  int        		`json:"valor"`
	Tipo string 			`json:"tipo"`
	Descricao string 		`json:"descricao"`
	RealizadaEm time.Time  	`json:"realizada_em"`	
}

type ExtratoController struct {
	idCliente 	int
	c 			*fiber.Ctx
}

func NewExtratoController(c *fiber.Ctx) (*ExtratoController) {
	return &ExtratoController{c:c}
}


func (e *ExtratoController) Run() ExtratoResult {

	idCliente, err := strconv.Atoi(e.c.Params("ID"))
	if err != nil  {
		return ExtratoResult{StatusCode: 422}
	}

	if idCliente < 1 || idCliente > 5 {
		return ExtratoResult{StatusCode: 404}
	}

	e.idCliente = idCliente

	result, err := e.run()
	if err != nil {
		return ExtratoResult{StatusCode: 501}
	}

	return *result
}


func (e *ExtratoController) run() (*ExtratoResult, error) {

	rows, err := helpers.GetDBConnection().Query(e.c.Context(),
		`	SELECT c.limite, c.saldo, t.valor, t.tipo, t.descricao, t.created_at
		    FROM clientes c LEFT JOIN transacoes t ON t.id_cliente = c.id
		    WHERE  c.id = $1
			ORDER BY t.created_at DESC
			LIMIT 10
		  `, 
		  e.idCliente,
	)
	defer rows.Close()

	if err != nil {
		logrus.Error("[ExtratoController] Erro ao obter extrato do cliente:", e.idCliente, ". Erro:", err)
		return nil, err
	}

	extratoResult := ExtratoResult{
		Saldo: ExtratoSaldoResult{},
		UltimasTransacoes: make([]ExtratoTransacaoResult, 0, 10),
	}

	for rows.Next() {
		var saldo int
		var limite int
		var tr ExtratoTransacaoResult

		err = rows.Scan(&limite, &saldo, &tr.Valor, &tr.Tipo, &tr.Descricao, &tr.RealizadaEm)
		if err != nil {
			extratoResult.UltimasTransacoes = make([]ExtratoTransacaoResult, 0, 0)
		} else {
			extratoResult.UltimasTransacoes = append(extratoResult.UltimasTransacoes, tr)
		}

		extratoResult.Saldo.Total = saldo
		extratoResult.Saldo.Limite = limite
		extratoResult.Saldo.DataExtrato = time.Now()	
	}

	extratoResult.StatusCode = 200

	return &extratoResult, nil
}
