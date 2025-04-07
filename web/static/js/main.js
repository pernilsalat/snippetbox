for (let link of document.querySelectorAll("nav a")) {
	if (link.getAttribute('href') === window.location.pathname) {
		link.classList.add("live");
		break;
	}
}