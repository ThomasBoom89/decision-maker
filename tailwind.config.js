/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: 'media',
    content: [
        "./views/**/*.html",
        "./frontend/**/*.js",
        "./internal/rendering/views/*.templ",
    ],
    theme: {
        extend: {},
    },
    plugins: [],
}

