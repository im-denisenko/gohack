# GoHack

[![Go](https://github.com/im-denisenko/gohack/actions/workflows/go.yml/badge.svg)](https://github.com/im-denisenko/gohack/actions/workflows/go.yml)

# Установка

```
curl -sLO https://github.com/im-denisenko/gohack/releases/download/v1.1.0/gohack-v1.1.0
chmod 755 gohack-v1.1.0
sudo mv gohack-v1.1.0 /usr/local/bin/gohack
```

# Использование

```
# Level 1
gohack -i input.json -o /dev/stdout -f json -a naive

# Level 2
gohack -i input.json -o output.json -f json -a naive

# Level 3
gohack -i input.json -o output.csv -f csv -a naive

# Level 4
gohack -i input.json -o output.db -f sqlite -a naive
```

# Флаги

- **`-h`** выводит все доступные флаги приложения.
- **`-a`** алгоритм генерации отчёта: naive, stream.
- **`-i`** путь к файлу, где взять список транзакций.
- **`-o`** путь к файлу, куда положить отчёт.
- **`-f`** формат отчёта: csv, json, sqlite.
- **`-q`** отключает вывод в консоль отладочных логов.

# Алгоритмы

## naive

Для расчётов полностью загружает входной файл в память.

## stream

Для расчётов использует потоковое чтение. В среднем медленнее, чем naive, но может обрабатывать файлы размером более, чем доступные RAM + SWAP.
