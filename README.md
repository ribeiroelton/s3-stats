# S3 STATS TOOL

A ferramenta S3 Stats Tool foi criada com o próposito de coletar informações de buckets S3 da AWS de uma determinada conta AWS. Foi desenvolvida utilizando Go 1.16 em conjunto com o SDK AWS GO V2. O processo de build foi realizado para 3 sistemas operacionais, sendo Windows, Linux e MacOS.

## Como Utilizar

A ferramenta funciona através da linha de comando e possui os binários disponibilizados no repo do Github

Releases <https://github.com/elribeiro/s3-stats-tool/releases/>

Os testes foram realizados em Sistema Operacional Linux e Windows

### No Linux

```bash
wget https://github.com/elribeiro/s3-stats-tool/releases/download/v0.0.3/s3analytics-linux-amd64

chmod +x s3analytics-linux-amd64
```

Para listar todas as opções:

```bash
./s3analytics-linux-amd64 -h 

Usage of ./s3analytics-linux-amd64:
  -fb string
    
                        String to filter only buckets that contains the specified value
                         (default no filter)
    
  -fo string
    
                        String to filter only objects that has a specific prefix
                         (default no filter) 
    
  -l
                        Boolean to define if this job will collect lifecycle rules as well
                         (default false)
    
  -o
                        Bool to indicate if output will be to a file named st3stats-date.json,
                        where date is the current date. If not set, will output to console in json format
                         (default false)
    
  -r
                        Boolean to define if this job will collect replication rules as well
                         (default false)
    
  -t int
    
                        Integer to define the number of threads to run concurrently.
                        Each thread will process one bucket at a time 
                        WATCH OUT: setting a high number may impact in high cost and computing resources usage
                         (default 2)
```

Para utilizar a ferramenta, é necessário realizar uma das duas opções de configuração de um profile AWS:

1. Exportar as variáveis de ambiente da AWS conforme exemplo no [link](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html). Essa é a configuração preferida visto que não precisa do client AWS instalado.

2. Configurar o profile "default" no client da AWS instalado no computador. Detalhes da configuração no [link](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html)

*IMPORTANTE:* O usuário utilizado para acessar a AWS deve possuir permissões de Leitura no S3 e seus recursos dependentes. Detalhes no [link](https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-policy-language-overview.html)

Exemplos de uso:

Buscar todos os buckets e seus respectivos objetos. Retorna para console em formato json

```bash
./s3analytics-linux-amd64
```

Buscar todos os buckets e seus respectivos objetos utilizando 10 threads de processamento. Retorna para console em formato json

```bash
./s3analytics-linux-amd64 -t 10
```

Buscar todos os buckets e seus respectivos objetos com informações de replicação e lifecycle. Retorna para console em formato json

```bash
./s3analytics-linux-amd64 -r -l
```

Busca informações apenas dos buckets que contenham ey7 no nome. Gera ouput para arquivo com nome padronizado em s3stats-AAAA-MM-DD.json

```bash
./s3analytics-linux-amd64 -o -fb ey7
```

## Como Contribuir

Esta ferramenta é sob licença MIT e para contribuir, basta forkar, gerar as alterações e enviar o PR :)

## Execução Local

O arquivo auxiliar Makefile (Utilizado no Linux), possui os atalhos para buildar, testar e executar a aplicação localmente

```bash

make build

make test

make run 

make all 
```

## Testes Comparativo Multithread / SingleThread

Cenário 1: 41 Buckets, sendo o maior com 20MM de arquivos.

```bash
time ./s3analytics-linux-amd64 -t 40

real	8m18.762s
user	3m10.360s
sys	0m17.585s
```

Cenário 2: 41 Buckets, sendo o maior com 20MM de arquivos.

```bash
time ./s3analytics-linux-amd64 -t 1

real	29m44.475s
user	3m31.435s
sys	0m33.447s
```
