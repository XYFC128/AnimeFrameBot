from src.main import start_bot, TMP_DIR
from telegram.ext import ApplicationBuilder
from unittest.mock import Mock, MagicMock
import configparser
import pytest


@pytest.fixture
def mock_os_makedirs(mocker):
    return mocker.patch('os.makedirs')


@pytest.fixture
def mock_config_parser(mocker):
    mock_config = {
        'Bot': {'name': '@MockBot', 'api_url': 'http://localhost:8763'},
        'Telegram Bot API': {'token': 'mock_token'},
    }
    mock_parser = MagicMock()
    mock_parser.__getitem__.side_effect = mock_config.__getitem__
    return mocker.patch.object(configparser.ConfigParser, '__new__', return_value=mock_parser)


@pytest.fixture
def mock_application_builder(mocker):
    mock_application = Mock()
    mock_application.token.return_value = mock_application
    mock_application.build.return_value = mock_application
    mocker.patch.object(mock_application, 'add_handler')
    mocker.patch.object(mock_application, 'run_polling')
    mocker.patch.object(ApplicationBuilder, '__new__', return_value=mock_application)
    return mock_application


def test_start_bot(
    mock_os_makedirs,
    mock_config_parser,
    mock_application_builder,
    mocker
):
    # Mock configuration path
    config_path = 'mock_config.ini'

    mocker.patch('telegram.ext.CommandHandler')
    mocker.patch('telegram.ext.MessageHandler')

    # Call start_bot
    start_bot(config_path)

    # Assertions
    mock_os_makedirs.assert_called_once_with(TMP_DIR, exist_ok=True)
    mock_config_parser.assert_called_once()
    
    # Assert polling was started
    mock_application_builder.run_polling.assert_called_once()
