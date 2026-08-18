#!/usr/bin/env bash
set -e

echo "Construindo imagem Docker para o Jogo de Damas..."
docker build -t damas-go:latest .

echo "Iniciando o jogo no terminal..."
docker run --rm -it damas-go:latest
