
export const updateTheme = (prefersDarkColorScheme: boolean) => {
    document.documentElement.classList.toggle('wa-dark', prefersDarkColorScheme);
}

export const watchColorSchemeChanges = (cb: (prefersDarkColorScheme: boolean) => void) => {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event =>
        cb(event.matches)
    )
};

export const prefersDarkColorScheme = () => {
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
}