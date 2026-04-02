package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sistema_venta_pasajes/internal/terminal/domain"
	"strconv"

	"sistema_venta_pasajes/internal/terminal/input"
	"sistema_venta_pasajes/internal/terminal/service"
	"sistema_venta_pasajes/internal/terminal/util"
	"sistema_venta_pasajes/pkg"

	"github.com/gorilla/mux"
)

type TerminalHandler struct {
	service service.TerminalService
}

func NewTerminalHandler(s service.TerminalService) *TerminalHandler {
	return &TerminalHandler{service: s}
}

func toTerminalOutput(t *domain.Terminal) input.TerminalOutput {
	return input.TerminalOutput{
		IDTerminal:   t.IDTerminal,
		Nombre:       t.NOMBRE,
		Ciudad:       t.CIUDAD,
		Departamento: t.DEPARTAMENTO,
		Direccion:    t.DIRECCION,
		Estado:       t.ESTADO,
	}
}

func (h *TerminalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateTerminalInput
	if err := decodeJSONBody(r, &in); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	terminal, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_TERMINAL_CREATED, toTerminalOutput(terminal), nil)
}

func (h *TerminalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := terminalIDFromRequest(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	terminal, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_TERMINAL_GET, toTerminalOutput(terminal), nil)
}

func (h *TerminalHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := terminalIDFromRequest(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	var in input.UpdateTerminalInput
	if err := decodeJSONBody(r, &in); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	terminal, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_TERMINAL_UPDATED, toTerminalOutput(terminal), nil)
}

func (h *TerminalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := terminalIDFromRequest(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_TERMINAL_DELETED, nil, nil)
}

func (h *TerminalHandler) List(w http.ResponseWriter, r *http.Request) {
	terminals, err := h.service.List()
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	outputs := make([]input.TerminalOutput, 0, len(terminals))
	for _, t := range terminals {
		outputs = append(outputs, toTerminalOutput(&t))
	}
	// Si no hay resultados, devolver array vacío
	var data any = outputs
	if len(outputs) == 0 {
		data = []interface{}{}
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_TERMINAL_LIST, data, map[string]any{"count": len(outputs)})
}

func terminalIDFromRequest(r *http.Request) (int64, error) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return 0, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_TERMINAL_INVALID_ID)
	}
	return id, nil
}

func decodeJSONBody(r *http.Request, destination any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return mapDecodeError(err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return pkg.BadRequest(util.ERR_CODE_INVALID_JSON, util.MSG_TERMINAL_SINGLE_JSON_OBJ).WithDetails(err.Error())
	}
	return nil
}

func mapDecodeError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		return pkg.BadRequest(util.ERR_CODE_EMPTY_BODY, "El cuerpo de la solicitud es obligatorio")
	case errors.As(err, &syntaxError):
		return pkg.BadRequest(util.ERR_CODE_INVALID_JSON, "El cuerpo JSON no tiene un formato valido")
	case errors.As(err, &unmarshalTypeError):
		return pkg.BadRequest(util.ERR_CODE_INVALID_JSON_TYPE, "Uno o mas campos del JSON tienen un tipo invalido")
	default:
		return pkg.BadRequest(util.ERR_CODE_INVALID_JSON, "El cuerpo JSON es invalido")
	}
}
