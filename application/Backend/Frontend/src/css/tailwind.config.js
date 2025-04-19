/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./application/Backend/Frontend/src/html/**/*.html",
    "./application/Backend/Frontend/src/views/**/*.hbs",
  ],
  theme: {
    extend: {
      colors: {
        'light-mode-bg': '#fffd98',
        'light-mode-header': '#BDE4A7',
      },
    },
  },
  plugins: [],
}

