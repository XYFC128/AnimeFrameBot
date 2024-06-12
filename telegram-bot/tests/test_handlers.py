from pytest_mock import MockerFixture
from src.main import *
from telegram import Update, InputMediaPhoto
from telegram.ext import ContextTypes
from unittest.mock import AsyncMock, Mock, mock_open
import pytest
import requests


@pytest.fixture
def mock_update(mocker: MockerFixture):
    mock_update = mocker.Mock(spec=Update)
    mock_update.effective_chat.id = 123456789
    mock_update.effective_chat.send_media_group = AsyncMock()
    return mock_update


@pytest.fixture
def mock_context(mocker: MockerFixture):
    mock_context = mocker.Mock(spec=ContextTypes.DEFAULT_TYPE)
    mock_context.bot.send_message = AsyncMock()
    return mock_context


@pytest.mark.asyncio
async def test_help_command(mock_update, mock_context):
    await help_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text=HELP_TEXT)


@pytest.mark.asyncio
async def test_start_command(mock_update, mock_context):
    await start(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text="I'm a bot, please talk to me!")


@pytest.mark.asyncio
async def test_frame_command(mock_update, mock_context, mocker: MockerFixture):
    # Mock context
    mock_context.args = ['query', '1']

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_response = mock_requests_get.return_value
    mock_response.json.return_value = [
        {'name': 'frame1', 'subtitle': 'subtitle1'},
    ]
    mock_response.content = b'image_content'

    # Mock InputMediaPhoto
    mock_input_media_photo = mocker.patch.object(InputMediaPhoto, '__new__', return_value=Mock())

    # Call the frame function
    await frame_command(mock_update, mock_context)

    # Assert query the api server
    mock_requests_get.assert_any_call(f"{API_URL}/frame/fuzzy/query/1")
    mock_requests_get.assert_any_call(f"{API_URL}/frame/frame1")

    # Assert that send_message was not called
    mock_context.bot.send_message.assert_not_called()

    # Assert that send_media_group was called with the correct parameters
    expected_media_group = [
        mock_input_media_photo(b'image_content', caption='subtitle1'),
    ]
    mock_update.effective_chat.send_media_group.assert_awaited_once_with(expected_media_group)

    mock_input_media_photo.assert_any_call(b'image_content', caption='subtitle1')


@pytest.mark.asyncio
async def test_frame_command_no_arg(mock_update, mock_context):
    mock_context.args = []
    await frame_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text="Please provide a text to frame like this: /frame {text}")


@pytest.mark.asyncio
async def test_frame_command_invalid_frame_num(mock_update, mock_context):
    mock_context.args = ['query', 'aaa']
    await frame_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a valid frame number.")

    mock_context.args = ['query', '-1']
    await frame_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a valid frame number.")

    mock_context.args = ['query', '000']
    await frame_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a valid frame number.")


@pytest.mark.asyncio
async def test_frame_command_frame_num_too_large(mock_update, mock_context):
    mock_context.args = ['query', '100']
    await frame_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text="I can only provide at most 10 frames at once.")


@pytest.mark.asyncio
async def test_frame_command_no_frame(mock_update, mock_context, mocker: MockerFixture):
    # Mock context
    mock_context.args = ['query', '1']

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_response = mock_requests_get.return_value
    mock_response.json.return_value = []

    # Call the frame function
    await frame_command(mock_update, mock_context)

    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Sorry, I'm unable to find any frame.")


@pytest.mark.asyncio
async def test_frame_command_request_error(mock_update, mock_context, mocker: MockerFixture):
    # Mock context
    mock_context.args = ['query', '1']

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_requests_get.side_effect = requests.exceptions.RequestException("Network error")

    # Call the frame function
    await frame_command(mock_update, mock_context)

    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Sorry, I'm unable to find any frame.")


@pytest.mark.asyncio
async def test_random_command(mock_update, mock_context, mocker: MockerFixture):
    # Mock context
    mock_context.args = ['1']

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_response = mock_requests_get.return_value
    mock_response.json.return_value = [
        {'name': 'frame1', 'subtitle': 'subtitle1'},
    ]
    mock_response.content = b'image_content'

    # Mock InputMediaPhoto
    mock_input_media_photo = mocker.patch.object(InputMediaPhoto, '__new__', return_value=Mock())

    # Call the frame function
    await random_command(mock_update, mock_context)

    # Assert query the api server
    mock_requests_get.assert_any_call(f"{API_URL}/frame/random/1")
    mock_requests_get.assert_any_call(f"{API_URL}/frame/frame1")

    # Assert that send_message was not called
    mock_context.bot.send_message.assert_not_called()

    # Assert that send_media_group was called with the correct parameters
    expected_media_group = [
        mock_input_media_photo(b'image_content', caption='subtitle1'),
    ]
    mock_update.effective_chat.send_media_group.assert_awaited_once_with(expected_media_group)

    mock_input_media_photo.assert_any_call(b'image_content', caption='subtitle1')


@pytest.mark.asyncio
async def test_random_command_invalid_frame_num(mock_update, mock_context):
    mock_context.args = ['aaa']
    await random_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a valid frame number.")

    mock_context.args = ['-1']
    await random_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a valid frame number.")

    mock_context.args = ['000']
    await random_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a valid frame number.")


@pytest.mark.asyncio
async def test_random_command_frame_num_too_large(mock_update, mock_context):
    mock_context.args = ['100']
    await random_command(mock_update, mock_context)
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text="I can only provide at most 10 frames at once.")


@pytest.mark.asyncio
async def test_random_command_no_frame(mock_update, mock_context, mocker: MockerFixture):
    # Mock context
    mock_context.args = ['1']

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_response = mock_requests_get.return_value
    mock_response.json.return_value = []

    # Call the frame function
    await random_command(mock_update, mock_context)

    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Sorry, I'm unable to find any frame.")


@pytest.mark.asyncio
async def test_random_command_request_error(mock_update, mock_context, mocker: MockerFixture):
    # Mock context
    mock_context.args = ['1']

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_requests_get.side_effect = requests.exceptions.RequestException("Network error")

    # Call the frame function
    await random_command(mock_update, mock_context)

    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Sorry, I'm unable to find any frame.")


@pytest.mark.asyncio
async def test_smart_reply_hit(mock_update, mock_context, mocker: MockerFixture):
    # Mock update
    mock_update.message.text = 'message'
    mock_update.message.id = 1234

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_response = mock_requests_get.return_value
    mock_response.json.return_value = [
        {'name': 'frame1', 'subtitle': 'subtitle1'},
    ]
    mock_response.content = b'image_content'

    # Mock InputMediaPhoto
    mock_input_media_photo = mocker.patch.object(InputMediaPhoto, '__new__', return_value=Mock())

    # Call the frame function
    await handle_smart_reply(mock_update, mock_context)

    # Assert query the api server
    mock_requests_get.assert_any_call(f"{API_URL}/frame/exact/message/1")
    mock_requests_get.assert_any_call(f"{API_URL}/frame/frame1")

    # Assert that send_message was not called
    mock_context.bot.send_message.assert_not_called()

    # Assert that send_media_group was called with the correct parameters
    expected_media_group = [
        mock_input_media_photo(b'image_content', caption='subtitle1'),
    ]
    mock_update.effective_chat.send_media_group.assert_awaited_once_with(expected_media_group, reply_to_message_id=mock_update.message.id)

    mock_input_media_photo.assert_any_call(b'image_content', caption='subtitle1')


@pytest.mark.asyncio
async def test_smart_reply_miss(mock_update, mock_context, mocker: MockerFixture):
    # Mock update
    mock_update.message.text = 'message'
    mock_update.message.id = 1234

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_response = mock_requests_get.return_value
    mock_response.json.return_value = [ ]

    # Call the frame function
    await handle_smart_reply(mock_update, mock_context)

    # Assert query the api server
    mock_requests_get.assert_any_call(f"{API_URL}/frame/exact/message/1")

    # Assert that no message was sent
    mock_context.bot.send_message.assert_not_called()
    mock_update.effective_chat.send_media_group.assert_not_called()


@pytest.mark.asyncio
async def test_smart_reply_error(mock_update, mock_context, mocker: MockerFixture):
    # Mock update
    mock_update.message.text = 'message'
    mock_update.message.id = 1234

    # Mock the responses for requests.get
    mock_requests_get = mocker.patch('requests.get')
    mock_requests_get.side_effect = requests.exceptions.RequestException("Network error")

    # Call the frame function
    await handle_smart_reply(mock_update, mock_context)

    # Assert query the api server
    mock_requests_get.assert_any_call(f"{API_URL}/frame/exact/message/1")

    # Assert that no message was sent
    mock_context.bot.send_message.assert_not_called()
    mock_update.effective_chat.send_media_group.assert_not_called()


@pytest.mark.asyncio
async def test_upload_photo_with_caption(mock_update, mock_context, mocker: MockerFixture):
    mock_update.message.caption = 'this is caption'

    mock_photo = Mock()
    mock_photo.file_id = 1234
    mock_update.message.photo = [mock_photo]

    mock_file = Mock()
    mock_file.download_to_drive = AsyncMock()
    mock_context.bot.get_file = AsyncMock()
    mock_context.bot.get_file.return_value = mock_file

    mock_upload = mocker.patch('src.main.upload', new=AsyncMock())

    await image_downloader(mock_update, mock_context)

    mock_context.bot.get_file.assert_awaited_once_with(mock_photo.file_id)
    mock_file.download_to_drive.assert_awaited_once()
    mock_upload.assert_awaited_once()


@pytest.mark.asyncio
async def test_upload_photo_without_caption(mock_update, mock_context):
    mock_update.message.caption = None

    mock_photo = Mock()
    mock_photo.file_id = 1234
    mock_update.message.photo = [mock_photo]

    await image_downloader(mock_update, mock_context)

    mock_context.bot.send_message.assert_awaited_with(chat_id=mock_update.effective_chat.id, text="Please provide a caption for the image")


@pytest.mark.asyncio
async def test_upload_file_with_caption(mock_update, mock_context, mocker: MockerFixture):
    mock_update.message.caption = 'this is caption'
    mock_update.message.document.file_id = 1234
    mock_update.message.document.file_name = 'filename.jpg'

    mock_file = Mock()
    mock_file.download_to_drive = AsyncMock()
    mock_context.bot.get_file = AsyncMock()
    mock_context.bot.get_file.return_value = mock_file

    mock_upload = mocker.patch('src.main.upload', new=AsyncMock())

    await image_file_downloader(mock_update, mock_context)

    mock_context.bot.get_file.assert_awaited_once_with(mock_update.message.document.file_id)
    mock_file.download_to_drive.assert_awaited_once()
    mock_upload.assert_awaited_once()


@pytest.mark.asyncio
async def test_upload_file_without_caption(mock_update, mock_context, mocker: MockerFixture):
    mock_update.message.caption = None
    mock_update.message.document.file_id = 1234
    mock_update.message.document.file_name = 'filename.jpg'

    mock_file = Mock()
    mock_file.download_to_drive = AsyncMock()
    mock_context.bot.get_file = AsyncMock()
    mock_context.bot.get_file.return_value = mock_file

    mock_upload = mocker.patch('src.main.upload', new=AsyncMock())

    await image_file_downloader(mock_update, mock_context)

    mock_context.bot.get_file.assert_awaited_once_with(mock_update.message.document.file_id)
    mock_file.download_to_drive.assert_awaited_once()
    mock_upload.assert_awaited_once()


@pytest.fixture
def mock_requests_post(mocker):
    return mocker.patch('requests.post')


@pytest.fixture
def mock_open_file(mocker):
    return mocker.patch('builtins.open', mock_open(read_data="data"))


@pytest.fixture
def mock_os_remove(mocker):
    return mocker.patch('os.remove')


@pytest.mark.asyncio
async def test_upload_success(mock_update, mock_context, mock_requests_post, mock_open_file, mock_os_remove):
    mock_response = Mock()
    mock_response.status_code = 201
    mock_requests_post.return_value = mock_response

    file_path = "/tmp/test_caption.jpg"

    await upload(mock_update, mock_context, file_path)

    # Assert that requests.post was called with the correct parameters
    mock_requests_post.assert_called_once_with(f"{API_URL}/frame", files={'image': mock_open_file.return_value})

    # Assert that send_message was called with the correct parameters
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text=f"Image {urllib.parse.unquote('test_caption.jpg')} uploaded successfully")

    # Assert that the file was closed and removed
    mock_open_file.return_value.close.assert_called_once()
    mock_os_remove.assert_called_once_with(file_path)


@pytest.mark.asyncio
async def test_upload_failure_status_code(mock_update, mock_context, mock_requests_post, mock_open_file, mock_os_remove):
    mock_response = Mock()
    mock_response.status_code = 400
    mock_response.text = "Bad request"
    mock_requests_post.return_value = mock_response

    file_path = "/tmp/test_caption.jpg"

    await upload(mock_update, mock_context, file_path)

    # Assert that requests.post was called with the correct parameters
    mock_requests_post.assert_called_once_with(f"{API_URL}/frame", files={'image': mock_open_file.return_value})

    # Assert that send_message was called with the error message
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text="Failed to upload due to an error in the api server.")

    # Assert that the file was closed and removed
    mock_open_file.return_value.close.assert_called_once()
    mock_os_remove.assert_called_once_with(file_path)


@pytest.mark.asyncio
async def test_upload_request_exception(mock_update, mock_context, mock_requests_post, mock_open_file, mock_os_remove):
    mock_requests_post.side_effect = requests.RequestException("Network error")

    file_path = "/tmp/test_caption.jpg"

    await upload(mock_update, mock_context, file_path)

    # Assert that requests.post was called with the correct parameters
    mock_requests_post.assert_called_once_with(f"{API_URL}/frame", files={'image': mock_open_file.return_value})

    # Assert that send_message was called with the error message
    mock_context.bot.send_message.assert_awaited_once_with(chat_id=mock_update.effective_chat.id, text="Failed to upload due to an internal error.")

    # Assert that the file was closed and removed
    mock_open_file.return_value.close.assert_called_once()
    mock_os_remove.assert_called_once_with(file_path)
