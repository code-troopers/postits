# Étape 1: Construction de l'application Go
FROM golang:1.23-alpine AS build-stage

# Définir le répertoire de travail
WORKDIR /app

# Copier le fichier go.mod et go.sum pour la gestion des dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code source
COPY . .

# Construire l'application
RUN go build -o main .

# Étape 2: Image exécutable
FROM alpine:latest

# Copier l'exécutable Go depuis le conteneur de build
COPY --from=build-stage /app/main /app/main

# Exposer le port 3000
EXPOSE 3000

# Démarrer l'application Go
CMD ["/app/main"]

