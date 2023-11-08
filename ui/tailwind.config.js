/** @type {import('tailwindcss').Config} */
const withMT = require("@material-tailwind/html/utils/withMT");

module.exports = withMT({
    content: ["./templates/*.templ"],
    plugins: [],
	theme: {
    extend: {},
  }
})
