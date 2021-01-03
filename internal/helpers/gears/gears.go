package gears

import (
	"fmt"
	"github.com/technopark_database/internal/models"
	"strings"
)

func AddPagination(query string, values []interface{},
	pagination *models.Pagination, since uint64, i int) (string, []interface{}) {
	if since != 0 {
		char := ""
		if pagination.Desc {
			char = "<"
		} else {
			char = ">"
		}
		createdString := "AND id" + char + fmt.Sprintf("$%d", i)
		query = strings.Join([]string{query,
			createdString,
		}, " ")
		i++
		values = append(values, since)
	}

	query = strings.Join([]string{query, "ORDER BY created"}, " ")
	if pagination.Desc {
		query = strings.Join([]string{query,
			"DESC",
		}, " ")
	}
	query = strings.Join([]string{query,
		", id",
	}, "")
	if pagination.Desc {
		query = strings.Join([]string{query,
			"DESC",
		}, " ")
	}

	limitStr := fmt.Sprintf("LIMIT $%d", i)
	query = strings.Join([]string{query, limitStr}, " ")
	values = append(values, pagination.Limit)

	//logrus.Info(query)
	//logrus.Info(values)

	return query, values
}
