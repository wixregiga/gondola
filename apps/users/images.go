package users

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"reflect"
	"strings"

	"gondola/app"
	"gondola/internal/httpserve"
	"gondola/net/httpclient"
)

var (
	ImageHandler app.Handler
	ImageFetcher func(ctx *app.Context, url string) (id string, format string, err error)

	imagePrefix string
)

func userImageId(ctx *app.Context, val reflect.Value) (string, string) {
	if image, _ := getUserValue(val, "Image").(string); image != "" {
		return image, getUserValue(val, "ImageFormat").(string)
	}
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", ""
		}
		val = val.Elem()
	}
	d := data(ctx)
	for _, v := range d.enabledSocialAccountTypes() {
		fval := val.FieldByName(v.Name.String())
		if fval.IsValid() && fval.Elem().IsValid() {
			image := fval.Elem().FieldByName("Image")
			if image.String() != "" {
				imageFormat := fval.Elem().FieldByName("ImageFormat")
				return image.String(), imageFormat.String()
			}
		}
	}
	return "", ""
}

func Image(ctx *app.Context, user interface{}) (string, error) {
	if imagePrefix == "" || user == nil {
		return "", nil
	}
	val := reflect.ValueOf(user)
	if !val.IsValid() {
		return "", nil
	}
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}
	if val.Type() != getUserType(ctx) {
		return "", fmt.Errorf("invalid user type %s, must be %s", val.Type(), getUserType(ctx))
	}
	id, format := userImageId(ctx, val)
	if id != "" {
		if format == "jpeg" {
			format = "jpg"
		}
		return imagePrefix + id + "." + format, nil
	}
	return "", nil
}

func UserImageHandler(ctx *app.Context) {
	if ImageHandler != nil {
		ImageHandler(ctx)
		return
	}
	id := ctx.IndexValue(0)
	format := ctx.IndexValue(1)
	if lower := strings.ToLower(format); lower != format {
		ctx.MustRedirectReverse(true, ImageHandlerName, id, lower)
		return
	}
	ctx.SetHeader("Content-Type", "image/"+format)
	httpserve.NeverExpires(ctx)
	bs := ctx.Blobstore()
	if err := bs.Serve(ctx, id, nil); err != nil {
		panic(err)
	}
}

func getImage(ctx *app.Context, url string) (string, string, error) {
	if ImageFetcher != nil {
		return ImageFetcher(ctx, url)
	}
	return defaultFetchImage(ctx, url)
}

func fetchImage(ctx *app.Context, url string) (string, string, string) {
	return mightFetchImage(ctx, url, "", "", "")
}

func mightFetchImage(ctx *app.Context, url string, prevId string, prevFormat string, prevURL string) (string, string, string) {
	if url == prevURL {
		return prevId, prevFormat, prevURL
	}
	if url == "" {
		if prevId != "" {
			ctx.Blobstore().Remove(prevId)
		}
		return "", "", ""
	}
	id, format, err := getImage(ctx, url)
	if err != nil {
		// Keep previous
		return prevId, prevFormat, prevURL
	}
	if prevId != "" {
		ctx.Blobstore().Remove(prevId)
	}
	return id, format, url
}

func defaultFetchImage(ctx *app.Context, url string) (string, string, error) {
	resp, err := httpclient.New(ctx).Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Close()
	data, err := resp.ReadAll()
	if err != nil {
		return "", "", err
	}
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", "", err
	}
	bs := ctx.Blobstore()
	id, err := bs.Store(data, nil)
	if err != nil {
		return "", "", err
	}
	return id, strings.ToLower(format), nil
}

func init() {
	app.Signals.WillListen.Listen(func(a *app.App) {
		placeholder := "0000placeholder0000"
		rev, err := a.Reverse(ImageHandlerName, placeholder, placeholder)
		if err == nil {
			if pos := strings.Index(rev, placeholder); pos >= 0 {
				imagePrefix = rev[:pos]
			}
		}
	})
}
