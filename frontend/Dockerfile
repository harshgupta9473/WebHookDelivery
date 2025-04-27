# Step 1: Build the Next.js app
FROM node:18-alpine AS build

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy all other files
COPY . .

# Build the Next.js app
RUN npm run build

# Step 2: Production image
FROM node:18-alpine AS production

WORKDIR /app

# Copy built app from the previous stage
COPY --from=build /app /app

# Install only production dependencies
RUN npm install --production

# Expose port 3000
EXPOSE 3000

# Start the app
CMD ["npm", "start"]
