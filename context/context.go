package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/labstack/echo"
	"gorm.io/gorm"

	vald "github.com/go-ozzo/ozzo-validation/v4"
)

type ContextIndex uint
type ParamError error

const (
	SessionContext ContextIndex = 10000
	ParamContext   ContextIndex = 10001

	originFormat = "2006-01-02" //yyyy-mm-dd
	paramFormat  = "02-01-2006" // dd-mm-yyy
)

var ErrQueryInvalid ParamError = errors.New("tipe data parameter tidak sesuai")

type Session interface {
	SetSession(interface{}) error
}

type Context struct {
	echo.Context
	session interface{}
	param   Param
}

type ContextModel interface {
	New(ctx *Context) ContextModel
	First(...interface{}) ContextModel
	Unscoped(...interface{}) ContextModel
	Where(string, ...interface{}) ContextModel
	Select(...string) ContextModel
	// Find(...interface{}) ContextModel
	Join(string, ...interface{}) ContextModel
	Limit(int) ContextModel
	Preload(string, ...interface{}) ContextModel
	Page(...interface{}) PageData
	Find(...interface{}) error
	Group(string) ContextModel
	Having(string, ...interface{}) ContextModel
	IsEmpty() bool
	HasError() bool
	GetError() error
	Filterizer
}

type PageData struct {
	Items          interface{} `json:"items"`
	Keyword        string      `json:"keyword"`
	ItemsPerPage   int         `json:"items_per_page"`
	TotalItems     int64       `json:"total_items"`
	Page           int         `json:"page"`
	TotalPage      int64       `json:"total_page"`
	Error          error       `json:"error,omitempty"`
	AcceptedFilter []string    `json:"accepted_filter,omitempty"`
}

func (p *PageData) HasError() bool {
	return p.Error != nil
}

func (p *PageData) GetError() error {
	return p.Error
}

type Range struct {
	RangedBy string `json:"ranged_by"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

func rangeRule() vald.Rule {
	return vald.By(func(value interface{}) (err error) {
		var (
			variable = value.(string)
			bytes    = []byte(variable)
		)
		switch {
		case len(variable) == 0:
		case !regexp.MustCompile(`^[a-zA-Z0-9_.]+$`).Match(bytes):
			err = errors.New("hanya menerima alfanumerik, underscore dan titik")
		case !regexp.MustCompile(`^[a-zA-Z]`).Match(bytes):
			err = errors.New("harus dimulai dengan alfabet")
		}
		return
	})
}

func startEndRule(opposite string) vald.Rule {
	return vald.By(func(value interface{}) (err error) {
		var variable = value.(string)

		if len(variable) == 0 {
			return nil
		}

		if len(opposite) == 0 {
			_, errDate := time.Parse(paramFormat, variable)
			_, errFloat := strconv.ParseFloat(variable, 64)

			if errFloat != nil && errDate != nil {
				return errors.New("hanya menerima tanggal (contoh: 31-12-2022) atau bilangan bulat/desimal (desimal menggunakan titik)")
			}
		} else {
			_, errVar := time.Parse(paramFormat, variable)
			_, errOps := time.Parse(paramFormat, opposite)

			if errVar == nil && errOps == nil {
				return nil
			}

			_, errVar = strconv.ParseFloat(variable, 64)
			_, errOps = strconv.ParseFloat(opposite, 64)

			if errVar == nil && errOps == nil {
				return nil
			}

			return errors.New("format start dan end tidak sama")
		}

		return nil
	})
}

func (form Range) Validate() error {
	return vald.ValidateStruct(&form,
		vald.Field(&form.RangedBy, rangeRule()),
		vald.Field(&form.Start, startEndRule(form.End)),
		vald.Field(&form.End, startEndRule(form.Start)),
	)
}

type Param struct {
	Keyword      string
	SearchBy     []string
	FilterBy     map[string]interface{}
	ItemsPerPage int
	Page         int
	OrderBy      string
	OrderMethod  string
	Range        Range
}

type ErrInvalidParam map[string]error

// Error returns the error string of Errors.
func (es ErrInvalidParam) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := make([]string, len(es))
	i := 0
	for key := range es {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	var s strings.Builder
	for i, key := range keys {
		if i > 0 {
			s.WriteString("; ")
		}
		if errs, ok := es[key].(ErrInvalidParam); ok {
			_, _ = fmt.Fprintf(&s, "%v: (%v)", key, errs)
		} else {
			_, _ = fmt.Fprintf(&s, "%v: %v", key, es[key].Error())
		}
	}
	s.WriteString(".")
	return s.String()
}

// MarshalJSON converts the Errors into a valid JSON.
func (es ErrInvalidParam) MarshalJSON() ([]byte, error) {
	errs := map[string]interface{}{}
	for key, err := range es {
		if ms, ok := err.(json.Marshaler); ok {
			errs[key] = ms
		} else {
			errs[key] = err.Error()
		}
	}
	return json.Marshal(errs)
}

type Filterizer interface {
	GetAcceptedFilter() []string
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func NewContext(c echo.Context) *Context {
	return c.(*Context)
}

// FormData parse the incoming POST body into "form" struct
// handle application/json and application/x-www-form-urlencoded
// and multipart/form-data
func (c *Context) FormData(form interface{}) error {
	var content string = c.Request().Header.Get("content-type")

	if strings.Contains(content, "application/json") {
		decoder := json.NewDecoder(c.Request().Body)
		decoder.Decode(form)
		return nil
	} else if strings.Contains(content, "application/x-www-form-urlencoded") {
		c.Request().ParseForm()
		decoder := schema.NewDecoder()
		decoder.Decode(form, c.Request().Form)
		return nil
	} else if strings.Contains(content, "multipart/form-data") {
		c.Request().ParseMultipartForm(1024 << 10)
		decoder := schema.NewDecoder()
		decoder.Decode(form, c.Request().Form)
		return nil
	}
	var err = vald.Errors{"error": errors.New("content-type data tidak sesuai")}
	return err
}

func (c *Context) FormValidatable(form vald.Validatable) (err error) {
	err = c.FormData(form)
	if err == nil {
		err = form.Validate()
	}
	return err
}

func (ctx *Context) GetParam() Param {
	return ctx.param
}

func (ctx *Context) SetParam(param Param) {
	ctx.param = param
}

func (ctx *Context) SetFilter(filter map[string]interface{}) {
	ctx.param.FilterBy = filter
}

func (ctx *Context) GetSession() interface{} {
	return ctx.session
}

func (ctx *Context) SetSession(session interface{}) (err error) {
	ctx.session = session
	if session == nil {
		err = errors.New("session tidak ditemukan")
	}
	return err
}

func MiddlewareChain(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var ctx = &Context{Context: c}

		page, _ := strconv.ParseInt(c.Request().URL.Query().Get("page"), 10, 64)
		limit, _ := strconv.ParseInt(c.Request().URL.Query().Get("items_per_page"), 10, 64)

		ranged := Range{
			RangedBy: c.Request().URL.Query().Get("ranged_by"),
			Start:    c.Request().URL.Query().Get("start"),
			End:      c.Request().URL.Query().Get("end"),
		}

		if err = ranged.Validate(); err != nil {
			return ctx.BadRequest(err)
		}

		var searchArr []string
		searchBy := strings.ReplaceAll(c.Request().URL.Query().Get("search_by"), " ", "")
		if len(searchBy) > 0 {
			searchArr = strings.Split(searchBy, ",")
		}

		ctx.SetParam(Param{
			Keyword:      c.Request().URL.Query().Get("keyword"),
			SearchBy:     searchArr,
			Page:         int(page),
			ItemsPerPage: int(limit),
			OrderBy:      c.Request().URL.Query().Get("order_by"),
			OrderMethod:  c.Request().URL.Query().Get("order_method"),
			Range: Range{
				RangedBy: c.Request().URL.Query().Get("ranged_by"),
				Start:    c.Request().URL.Query().Get("start"),
				End:      c.Request().URL.Query().Get("end"),
			},
		})

		if ctx.Request().Context().Value(SessionContext) != nil {
			if err = ctx.SetSession(
				ctx.Request().Context().Value(SessionContext)); err != nil {
				return ctx.ServerError(err)
			}
		}

		return next(ctx)
	}
}

type Tabler interface {
	TableName() string
}

// Pagination set offset and limit of query.
// Default value of limit is 10, and offset of page 1.
// Note: dateColumn is an optional parameter and the function
// only use dateColumn[0]. Call this function before PageResult
func Paginate(_db *gorm.DB, query Param, dest interface{}) PageData {
	var (
		withKeyword bool = false
		keyword     string
		limit       int = 10
		page        int = 0
		between     string
		gs          *gorm.DB        = _db.Session(&gorm.Session{NewDB: true})
		stm         *gorm.Statement = &gorm.Statement{DB: _db}
	)

	if query.ItemsPerPage > 0 {
		limit = query.ItemsPerPage
	}

	stm.Parse(dest)

	if len(query.SearchBy) > 0 && len(query.Keyword) > 0 {
		for _, column := range query.SearchBy {
			keyword = fmt.Sprintf("cast(%s as varchar) ilike ?", column)
			gs = gs.Or(keyword, fmt.Sprintf("%%%s%%", query.Keyword))
			withKeyword = true
		}
	} else if len(query.Keyword) > 0 {
		result, _ := gs.Migrator().ColumnTypes(dest)
		for _, v := range result {
			keyword = fmt.Sprintf("cast(%s.%s as varchar) ilike ?", stm.Schema.Table, v.Name())
			gs = gs.Or(keyword, fmt.Sprintf("%%%s%%", query.Keyword))
			withKeyword = true
		}
	}

	if withKeyword {
		_db = _db.Where(gs)
	}

	_db = _db.Limit(limit)
	page = (query.Page - 1) * limit
	_db = _db.Offset(page)

	if query.Range.RangedBy != "" {

		startDate, errStart := time.Parse(paramFormat, query.Range.Start)
		endDate, errEnd := time.Parse(paramFormat, query.Range.End)

		switch {
		case len(query.Range.Start) > 0 && len(query.Range.End) > 0:
			between = fmt.Sprintf("date(%s) between ? and ?", query.Range.RangedBy)
			if errStart == nil && errEnd == nil {
				_db = _db.Where(between, startDate.Format(originFormat), endDate.Format(originFormat))
			} else {
				between = fmt.Sprintf("cast(%s as double precision) between cast(? as double precision) and cast(? as double precision)", query.Range.RangedBy)
				_db = _db.Where(between, query.Range.Start, query.Range.End)
			}
		case len(query.Range.Start) > 0:
			between = fmt.Sprintf("date(%s) >= ?", query.Range.RangedBy)
			if errStart == nil {
				_db = _db.Where(between, startDate)
			} else {
				between = fmt.Sprintf("cast(%s as double precision) >= cast(? as double precision)", query.Range.RangedBy)
				_db = _db.Where(between, query.Range.Start)
			}
		case len(query.Range.End) > 0:
			between = fmt.Sprintf("date(%s) <= ?", query.Range.RangedBy)
			if errEnd == nil {
				_db = _db.Where(between, endDate)
			} else {
				between = fmt.Sprintf("cast(%s as double precision) <= cast(? as double precision)", query.Range.RangedBy)
				_db = _db.Where(between, query.Range.End)
			}
		}
	}

	_db = _db.Find(dest)
	if _db.Error != nil {
		return PageData{
			Error:        _db.Error,
			Items:        []string{},
			Keyword:      query.Keyword,
			ItemsPerPage: limit,
			TotalPage:    1,
			Page:         1,
		}
	}

	var result = paginate(_db, dest, query)
	result.Error = _db.Error
	return result
}

// PageResult handle detail of BuildQuery and returns PaginationResult.
// This should called after Pagination() function to generate offset limit.
// Warning: must supply executed gorm.DB object () ex: PageResult(db.Find(&users), urlQuery)
func paginate(_db *gorm.DB, dest interface{}, query Param) PageData {
	var (
		count int64
		limit int = 10
		pager int = 0
		page  PageData
	)

	_db.Offset(-1).Limit(-1).Count(&count)

	if len(query.OrderBy) > 0 {
		_db = _db.Group(query.OrderBy)
	}

	if query.ItemsPerPage > 0 {
		limit = query.ItemsPerPage
		if limit > 100 {
			limit = 100
		}
	}

	if query.Page == 0 {
		query.Page = 1
	}

	pager = (query.Page - 1) * limit
	page.Items = dest
	page.Keyword = query.Keyword
	page.ItemsPerPage = limit
	page.Page = pager/limit + 1
	page.TotalItems = count
	page.TotalPage = int64(math.Ceil(float64(count) / float64(limit)))

	if page.TotalPage == 0 {
		page.TotalPage = 1
	}

	return page
}

func (ctx *Context) Success(payload interface{}) error {

	return ctx.JSON(http.StatusOK,
		Response{
			Status: "berhasil",
			Data:   payload,
		})
}

func (ctx *Context) NotFound(err error) error {

	return ctx.JSON(http.StatusNotFound,
		Response{
			Status: "record tidak ditemukan",
			Data: vald.Errors{
				"message": err,
			},
		})
}

func (ctx *Context) BadRequest(err error) error {
	return ctx.JSON(http.StatusBadRequest,
		Response{
			Status: "input belum sesuai",
			Data: vald.Errors{
				"message": err,
			},
		})
}

func (ctx *Context) Conflict(err error) error {
	return ctx.JSON(
		http.StatusBadRequest,
		Response{
			Status: "data duplikat",
			Data: vald.Errors{
				"message": err,
			},
		})
}

func (ctx *Context) Unauthorized(err error) error {
	return ctx.JSON(
		http.StatusUnauthorized,
		Response{
			Status: "user tidak teridentifikasi",
			Data: vald.Errors{
				"message": err,
			},
		})
}

func (ctx *Context) Forbidden(err error) error {
	return ctx.JSON(
		http.StatusForbidden,
		Response{
			Status: "akses dibatasi",
			Data: vald.Errors{
				"message": err,
			},
		},
	)
}

func (ctx *Context) ServerError(err error) error {
	var response Response

	response.Status = "internal server error"

	_, invalidParam := err.(ErrInvalidParam)

	if invalidParam {
		response.Status = "parameter tidak sesuai"
		response.Data = err
		return ctx.JSON(http.StatusBadRequest, response)
	}

	switch os.Getenv("DEBUG") {
	case "1", "true":
		response.Data = vald.Errors{"message": err}
	default:
		response.Data = vald.Errors{"message": errors.New("terjadi kendala teknis")}
	}

	return ctx.JSON(http.StatusInternalServerError, response)
}
