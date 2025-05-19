/**
 * @file: errors.go
 * @description: Базовые ошибки для BloFin REST-клиента
 * @dependencies: -
 * @created: 2025-05-19
 */

package utils

import "errors"

var ErrAPIRequest = errors.New("blofin: api request error")
