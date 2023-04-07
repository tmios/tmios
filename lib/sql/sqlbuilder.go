package sql

import (
	"fmt"
	"gorm.io/gorm"
)

type WhereExpr struct {
	LikeLR bool
	Expr   string
	Args   []string
}

type ReverseWhereExpr struct {
	Expr interface{}
	Args []interface{}
}

type JoinExpor struct {
	Expr string
	Args []interface{}
}

type sqlBuilder struct {
	whereExprs        []WhereExpr
	reverseWhereExprs []ReverseWhereExpr
	preloads          []string
	orderExprs        []interface{}
	joinExprs         []JoinExpor
	whereFuncs        []WhereFunc
}

func (s *sqlBuilder) Where(expr string, args ...string) *sqlBuilder {
	s.whereExprs = append(s.whereExprs, WhereExpr{false, expr, args})
	return s
}

func (s *sqlBuilder) LikeLR(expr string, args ...string) *sqlBuilder {
	s.whereExprs = append(s.whereExprs, WhereExpr{true, expr, args})
	return s
}

func (s *sqlBuilder) ReverseWhere(expr interface{}, args ...interface{}) *sqlBuilder {
	s.reverseWhereExprs = append(s.reverseWhereExprs, ReverseWhereExpr{expr, args})
	return s
}

func (s *sqlBuilder) Joins(expr string, args ...interface{}) *sqlBuilder {
	s.joinExprs = append(s.joinExprs, JoinExpor{Expr: expr, Args: args})
	return s
}

func (s *sqlBuilder) Preload(expr string) *sqlBuilder {
	s.preloads = append(s.preloads, expr)
	return s
}

func (s *sqlBuilder) Order(expr interface{}) *sqlBuilder {
	s.orderExprs = append(s.orderExprs, expr)
	return s
}

func (s *sqlBuilder) WhereFunc(whereFunc WhereFunc) *sqlBuilder {
	s.whereFuncs = append(s.whereFuncs, whereFunc)
	return s
}

func (s *sqlBuilder) Clone() *sqlBuilder {
	newS := *s
	return &newS
}

func (s *sqlBuilder) Build(query map[string]interface{}) WhereFunc {
	return func(q *gorm.DB) *gorm.DB {
		for _, joinExpr := range s.joinExprs {
			q = q.Joins(joinExpr.Expr, joinExpr.Args...)
		}

		for _, whereExpr := range s.whereExprs {
			var queryVals []interface{}
			for _, arg := range whereExpr.Args {
				val, ok := query[arg]
				if ok {
					if whereExpr.LikeLR {
						queryVals = append(queryVals, fmt.Sprintf("%v%v%v", "%", val, "%"))
					} else {
						queryVals = append(queryVals, val)
					}

				}
			}
			if len(queryVals) == len(whereExpr.Args) {
				q = q.Where(whereExpr.Expr, queryVals...)
			}
		}

		for _, expr := range s.reverseWhereExprs {
			q = q.Where(expr.Expr, expr.Args...)
		}

		for _, preload := range s.preloads {
			q = q.Preload(preload)
		}

		for _, orderExpr := range s.orderExprs {
			q = q.Order(orderExpr)
		}

		for _, whereFunc := range s.whereFuncs {
			q = whereFunc(q)
		}

		return q
	}
}

func Builder() *sqlBuilder {
	return &sqlBuilder{}
}
