from src.main import escape_path

def test_escape_path():
    assert escape_path('aaa') == 'aaa'
    assert escape_path('中文') == '中文'
    assert escape_path('.....data') == 'data'
    assert escape_path('a/b/c') == 'a_b_c'
    assert escape_path('.././') == '_._'
