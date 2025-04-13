package utils

import "golang.org/x/sys/windows"

func GetFileID(path string) (*uint64, error) {
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var info windows.ByHandleFileInformation

	err = windows.GetFileInformationByHandle(handle, &info)
	if err != nil {
		return nil, err
	}

	fileID := (uint64(info.FileIndexHigh) << 32) | uint64(info.FileIndexLow)

	return &fileID, nil
}
