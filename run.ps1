Write-Host "Construindo imagem Docker para o Jogo de Damas..." -ForegroundColor Cyan
docker build -t damas-go:latest .

if ($LASTEXITCODE -eq 0) {
    Write-Host "Iniciando o jogo no terminal..." -ForegroundColor Green
    docker run --rm -it damas-go:latest
} else {
    Write-Host "Falha na construcao da imagem Docker." -ForegroundColor Red
}
