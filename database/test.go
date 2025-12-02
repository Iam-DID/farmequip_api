package database

import (
    "context"
    "errors"
    "strings"

    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func Bool(v bool) *bool {
    return &v
}

func CheckCloudinaryConnection(cld *cloudinary.Cloudinary) error {

    _, err := cld.Upload.Upload(
        context.Background(),
        strings.NewReader("ping-test"),
        uploader.UploadParams{
            Folder:        "connection_test",
            PublicID:      "ping_check",
            Overwrite:     Bool(false),
            UniqueFilename: Bool(false),
            UseFilename:    Bool(true),
        },
    )

    if err != nil {
        return errors.New("Cloudinary connection FAILED: " + err.Error())
    }

    return nil
}
