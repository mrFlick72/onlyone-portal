FROM node:22-alpine3.20

EXPOSE 6060 6061

COPY dist .
COPY node_modules node_modules
CMD ["node", "app/index.js"]