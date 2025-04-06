# Use the official Node.js image for building the React app
FROM node:23 as builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Expose the frontend port
EXPOSE 80

# Start Nginx
CMD ["npm","run", "dev"]