import { defineConfig, loadEnv, type Plugin } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import fs from 'fs';

const pages: Record<string, { html: string; entry: string }> = {
  home:          { html: 'index.html',                    entry: '/home/index.tsx' },
  callback:      { html: 'callback.html',                 entry: '/auth/Callback.tsx' },
  logout:        { html: 'logout.html',                   entry: '/auth/Logout.tsx' },
  budgetExpense: { html: 'budget/expense/index.html',     entry: '/budget/index.tsx' },
  budgetRevenue: { html: 'budget/revenue/index.html',     entry: '/budget/index.tsx' },
  budgetTags:    { html: 'budget/search-tags/index.html', entry: '/budget/index.tsx' },
  account:       { html: 'account/index.html',            entry: '/account/index.tsx' },
};

function htmlTemplatePlugin(): Plugin {
  const template = fs.readFileSync(path.resolve(__dirname, 'template.html'), 'utf-8');

  for (const { html: htmlPath, entry } of Object.values(pages)) {
    const fullPath = path.resolve(__dirname, htmlPath);
    fs.mkdirSync(path.dirname(fullPath), { recursive: true });
    fs.writeFileSync(
      fullPath,
      template.replace('</body>', `  <script type="module" src="${entry}"></script>\n</body>`),
    );
  }

  return {
    name: 'html-template',
    closeBundle() {
      for (const { html: htmlPath } of Object.values(pages)) {
        console.log("closeBundle called")
        try { fs.unlinkSync(path.resolve(__dirname, htmlPath)); } catch {}
      }
    },
  };
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(__dirname, '../environments'), '');

  return {
    root: __dirname,
    plugins: [htmlTemplatePlugin(), react()],
    define: {
      'import.meta.env.REDIRECT_URI':                  JSON.stringify(env.REDIRECT_URI),
      'import.meta.env.CLIENT_APPLICATION_ID':         JSON.stringify(env.CLIENT_APPLICATION_ID),
      'import.meta.env.IDP_BASE_URL':                  JSON.stringify(env.IDP_BASE_URL),
      'import.meta.env.AUTHENTICATION_CHECK_INTERVAL': JSON.stringify(env.AUTHENTICATION_CHECK_INTERVAL),
      'import.meta.env.BUDGET_API_BASE_URL':           JSON.stringify(env.BUDGET_API_BASE_URL),
      'import.meta.env.REVENUE_API_BASE_URL':          JSON.stringify(env.REVENUE_API_BASE_URL),
      'import.meta.env.ACCOUNT_API_BASE_URL':          JSON.stringify(env.ACCOUNT_API_BASE_URL),
      'import.meta.env.TAG_API_BASE_URL':              JSON.stringify(env.TAG_API_BASE_URL),
    },
    build: {
      outDir: path.resolve(__dirname, '../dist'),
      emptyOutDir: true,
      rollupOptions: {
        input: Object.fromEntries(
          Object.entries(pages).map(([name, { html }]) => [
            name,
            path.resolve(__dirname, html),
          ]),
        ),
      },
    },
    server: { port: 3000 },
  };
});
