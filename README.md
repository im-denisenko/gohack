# GoHack

[![Go](https://github.com/im-denisenko/gohack/actions/workflows/go.yml/badge.svg)](https://github.com/im-denisenko/gohack/actions/workflows/go.yml)

# Установка

```
curl -sLO https://github.com/im-denisenko/gohack/releases/download/v1.0.0/gohack-v1.0.0
chmod 755 gohack-v1.0.0
sudo mv gohack-v1.0.0 /usr/local/bin/gohack
```

# Использование

```
gohack -i input.json -o output.json -f json -a naive
```

# Флаги

- **-h** выводит все доступные флаги приложения.
- **-a** алгоритм генерации отчёта: naive, stream.
- **-i** путь к файлу, где взять список транзакций.
- **-o** путь к файлу, куда положить отчёт.
- **-f** формат отчёта: csv, json.

# Алгоритмы

## naive

Для расчётов полностью загружает входной файл в память.

## stream

Для расчётов использует потоковое чтение. В среднем медленнее, чем naive, но может обрабатывать файлы размером более, чем доступные RAM + SWAP.
