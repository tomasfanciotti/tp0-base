
from .utils import Bet, store_bets


def register_bet(argv):
    new_bet = Bet(argv[0], argv[1], argv[2], argv[3], argv[4], argv[5])
    store_bets([new_bet])

    return argv[3], argv[5]
