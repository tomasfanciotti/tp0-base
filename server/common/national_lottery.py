from .utils import Bet, store_bets,load_bets,has_won

AGENCY_EXPECTED = 5
agency_readyness = []

def register_bet(argv):
    new_bet = Bet(argv[0], argv[1], argv[2], argv[3], argv[4], argv[5])
    store_bets([new_bet])

    return argv[3], argv[5]


def register_batch(argv):
    total_bets = int(argv.pop(0))

    bet_nro = 0
    buffer = []

    while len(buffer) < total_bets:

        idx = bet_nro * 7
        actual = argv[idx]
        if actual == f"@{bet_nro}":
            try:
                new_bet = Bet(argv[idx + 1], argv[idx + 2], argv[idx + 3], argv[idx + 4], argv[idx + 5], argv[idx + 6])
                buffer.append(new_bet)
                bet_nro += 1
            except Exception as e:
                print(e)
                return 0
        else:
            print("No cumple la condición de control")
            return 0

    store_bets(buffer)

    return len(buffer)


def ready(argv):

    agency_id = argv.pop(0)
    agency_readyness.append(agency_id)

    if len(agency_readyness) != AGENCY_EXPECTED:
        return agency_id, False

    return agency_id, True


def ask_winner(argv):

    agency_id = argv.pop(0)

    if len(agency_readyness) != AGENCY_EXPECTED:
        return agency_id, None
    
    winners = []
    bets = load_bets()
    for bet in bets:
        if bet.agency == int(agency_id) and has_won(bet):
            winners.append(bet.document)

    return agency_id, winners
