# Realtime Production Debugging

A high-performance, interactive presentation built with React and Tailwind CSS, focusing on **Mechanical Sympathy** and production profiling in Go.

## 🚀 Running Locally

This project uses **Native ES Modules** and **Import Maps**, meaning it requires no complex build step (like Webpack or Babel) to run for development.

### Prerequisites
- [Node.js](https://nodejs.org/) (for serving the files)

### Steps
1. **Clone/Download** the project directory.
2. Open your terminal in the project root.
3. Start a local development server:
   ```bash
   # Using npx (easiest)
   npx serve .

   # OR using Python
   python3 -m http.server 8000
   ```
4. Open your browser to the provided URL (usually `http://localhost:3000` or `http://localhost:8000`).

---

## 🛠 Project Structure

- `index.html`: The entry point. It contains the Tailwind configuration and the browser's `importmap` for React.
- `index.tsx`: Initializing the React root.
- `App.tsx`: The "Slide Deck" controller. Contains the slide data array and navigation logic.
- `components/`: Modular slide elements (Terminal, Gopher, GrafanaGraph, etc.).
- `constants.ts`: Branding colors and shared UI constants.

---

## 📦 Generating a Static Version

Because this application is built using standard web technologies and browser-native modules, it is **already a static application**.

### 1. Simple Deployment (Recommended)
You can host this entire folder on any static hosting provider (GitHub Pages, Netlify, Vercel, or an S3 bucket) exactly as it is. There is no `npm run build` required because the browser handles the module loading via the `importmap` defined in `index.html`.

### 2. "Next.js Style" Bundling (For Optimization)
If you want a highly optimized, minified, single-bundle version (similar to `next export`), you can use **Vite**:

1. Initialize a Vite project: `npm create vite@latest . --template react-ts`
2. Move the files into the structure Vite expects.
3. Run `npm run build`.
4. This will generate a `dist/` folder containing a traditional minified static site.

### 3. PDF/Offline Export
To get a static "snapshot" for sharing:
- Open the app in Chrome.
- Use `Cmd + P` (or `Ctrl + P`) and select **Save as PDF**.
- Ensure "Background Graphics" is checked in the print settings to keep the GoFr dark-mode aesthetics.

---

## ⌨️ Navigation Shortcuts
- `Space` / `Right Arrow`: Next Slide
- `Left Arrow`: Previous Slide
- `F`: Toggle Fullscreen
- `S`: Toggle Speaker Notes
