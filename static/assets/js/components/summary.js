PetiteVue.createApp({
    $delimiters: ['${', '}'],
    activityChartSvg: '',
    get currentInterval() {
        const urlParams = new URLSearchParams(window.location.search)
        if (urlParams.has('interval')) return urlParams.get('interval')
        if (!urlParams.has('from') && !urlParams.has('to')) return 'today'
        return null
    },
    mounted({ userId }) {
        const isDarkMode = window.matchMedia(
            '(prefers-color-scheme: dark)'
        ).matches
        const darkParam = isDarkMode ? 'dark' : ''
        fetch(`api/activity/chart/${userId}.svg?${darkParam}&noattr`)
            .then((res) => res.text())
            .then((data) => (this.activityChartSvg = data))
    },
}).mount('#summary-page')
