# Étape 1: Construire l'application Vue.js
FROM node:18 AS build-stage

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers de configuration
COPY package*.json ./

# Installer les dépendances
RUN npm install

# Copier le reste du code source
COPY . .

# Construire l'application pour la production
RUN npm run build

# Étape 2: Servir l'application avec Nginx
FROM nginx:stable-alpine AS production-stage

# Copier les fichiers construits dans le répertoire Nginx
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copier le fichier de configuration Nginx par défaut
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Exposer le port 80
EXPOSE 5500

# Démarrer Nginx
CMD ["nginx", "-g", "daemon off;"]
