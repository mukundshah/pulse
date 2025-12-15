export const decomposeURL = (url: string) => {
  const urlObj = new URL(url)

  const secure = urlObj.protocol === 'https:'

  return {
    host: urlObj.hostname,
    port: urlObj.port ? Number.parseInt(urlObj.port) : (secure ? 443 : 80),
    path: urlObj.pathname,
    queryParams: Object.fromEntries(urlObj.searchParams.entries()),
    secure,
  }
}

export const constructURL = ({ host, port, path, queryParams, secure }: { host: string, port?: number, path?: string, queryParams?: Record<string, string>, secure?: boolean }) => {
  const protocol = secure ? 'https:' : 'http:'

  const url = new URL(`${protocol}//${host}`)

  // only include port if it's non-default
  if (
    port
    && !(
      (secure && port === 443)
      || (!secure && port === 80)
    )
  ) {
    url.port = String(port)
  }

  url.pathname = path ?? ''
  url.search = new URLSearchParams(queryParams ?? {}).toString()

  return url.toString()
}
