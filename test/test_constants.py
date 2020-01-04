"""Test if the constant values are set up correctly."""
import src.utils.constants as constants


def test_constants_set_correctly():
    assert constants.SECOND == 1, 'Second should be the smallest unit!'
    assert constants.MINUTE == 60, 'Minute should be defined as seconds!'
    assert constants.HOUR == 3600, 'Hour should be defined as seconds!'
    assert constants.DAY == 86400, 'Day should be defined as seconds!'
    assert constants.WEEK == 604800, 'Week should be defined as seconds!'
    assert constants.MONTH == 2592000, 'Month should be defined as seconds!'
    assert constants.YEAR == 31536000, 'Year should be defined as seconds!'
    assert constants.CONVERT_PERCENT == 100, '''Second should be the smallest
unit!'''
    assert constants.CONVERT_KILO == 1e3, '''KILO means multiply by 1000!'''
    assert constants.CONVERT_NANO == 1e-9, '''KILO means multiply by 1000!'''
    assert constants.CONVERT_GIGA == 1e9, '''KILO means multiply by 1000!'''
