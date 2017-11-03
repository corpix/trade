#!/usr/bin/env python3
from http.client import HTTPSConnection
from re import findall, match
from sys import stderr, stdout


class Request(object):
    headers = {
        "Accept": "text/html",
        "Accept-Charset": "utf-8",
        "Cache-Control": "no-cache",
    }

    def __init__(self, host, timeout=15):
        self.host = host
        self.timeout = timeout

    def send(self, url, method, body, headers={}):
        request_headers = Request.headers.copy()
        request_headers.update(headers)

        connection = HTTPSConnection(
            self.host,
            timeout=self.timeout
        )
        connection.request(
                method,
                url,
                body=body,
                headers=request_headers
            )
        res = connection.getresponse()

        return res

def crypto_currencies():
    r = Request(host="coinmarketcap.com")
    response = r.send(
        "/all/views/all/",
        "GET",
        None
    )

    text = "".join(
        response.read().decode("utf-8").split("\n")
    )
    response.close()

    table = "".join(
        findall(
            """<table.*?id="currencies-all".*?</table>""",
            text
        )
    )
    body = "".join(
        findall(
            "<tbody>.*?</tbody>",
            table
        )
    )
    currencies = findall(
        "<tr.*?</tr>",
        body
    )
    for currency in currencies:
        columns = findall(
            "<td.*?</td>",
            currency
        )
        if match(".*?low vol.*?", columns[7].lower()):
            continue
        yield match(
            (
                ".*?<td.*?>.*?<a.*?href=\"/currencies/(.*?)\".*?</td>.*?"
                ".*?<td.*?>(.*?)</td>.*?"
            ),
            columns[1] + columns[2]
        ).groups()

def fiat_currencies():
    return [
        ("China Yan", "CNY"),
        ("Japanese Yen", "JPY"),
        ("Russian Ruble", "RUB"),
        ("United States Dollar", "USD"),
        ("Euro", "EUR"),
        ("Canadian Dollar", "CAD"),
    ]


def format_name(s):
    alphanumeric = list(range(ord("a"), ord("z")+1)) + \
                   list(range(ord("A"), ord("Z")+1)) + \
                   list(range(ord("0"), ord("9")+1))

    return "".join(
        [
            "Coin" + v.capitalize() if v[0].isnumeric() else v.capitalize()
            for v in "".join(
                    map(
                        lambda v: v if ord(v) in alphanumeric else " ",
                        s.strip()
                    )
            ).split(" ")
            if len(v) > 0
        ]
    )

def format_currencies(currencies):
    res = {}
    for name, code in sorted(currencies):
        ident = format_name(name)
        if ident in res:
            stderr.write(
                "Ident '{}' with code '{}' already in currencies with code '{}', not rewriting\n".format(
                    ident,
                    code,
                    res[ident]
                )
            )
            continue
        res[ident] = code.upper()
    return res

def generate_code(currencies):
    return """package currencies

const (
	InvalidCurrency Currency = iota
	""" + "\n\t".join(
            sorted(currencies.keys())
        ) + """
)

var (
	CurrencyMapping = map[Currency]string{
		""" + "\n\t\t".join(
                    [
                        "{}: \"{}\",".format(k, v)
                        for k, v in sorted(currencies.items())
                    ]
                ) + """
	}
)
"""


def main():
    stdout.write(
        generate_code(
            format_currencies(
                list(
                    crypto_currencies()
                ) + list(
                    fiat_currencies()
                )
            )
        ) + "\n"
    )

if __name__ == "__main__":
    main()
