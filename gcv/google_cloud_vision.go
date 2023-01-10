package gcv

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/vision/v1"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

func DefaultDialer(file_dir, type_req string) *vision.BatchAnnotateImagesResponse {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, vision.CloudVisionScope)
	if err != nil {
		log.Println(err)
	}

	service, err := vision.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println(err)
	}

	b, err := ioutil.ReadFile(file_dir)
	if err != nil {
		log.Println(err)
	}

	req := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},

		Features: []*vision.Feature{
			{
				Type:       type_req,
				MaxResults: 50,
			},
		},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		log.Println(err)
	}

	return res
}

//label-detection

func LabelOcr(file_dir string) string {
	var label_annotations []string
	res := DefaultDialer(file_dir, "LABEL_DETECTION")

	if annotations := res.Responses[0].LabelAnnotations; len(annotations) > 0 {
		for _, a := range annotations {
			label := a.Description
			score := strconv.Itoa(int(a.Score * 100))

			label_annotations = append(label_annotations, "{`Label`: "+label+", `Score`: "+score+"}")
		}
	}

	log.Println("label", strings.Join(label_annotations[:], "."))
	return strings.Join(label_annotations[:], ".")
}

//logo-detection

func LogoOcr(file_dir string) string {
	var logo_annotations []string
	res := DefaultDialer(file_dir, "LOGO_DETECTION")

	if annotations := res.Responses[0].LogoAnnotations; len(annotations) > 0 {
		for _, a := range annotations {
			logo := a.Description
			score := strconv.Itoa(int(a.Score * 100))

			logo_annotations = append(logo_annotations, "{`Logo`: "+logo+", `Score`: "+score+"}")
		}
	}

	log.Println("logo", strings.Join(logo_annotations[:], "."))
	return strings.Join(logo_annotations[:], ".")
}

//face-detection

func FaceOcr(file_dir string) string {
	var face_annotations []string
	res := DefaultDialer(file_dir, "FACE_DETECTION")

	if annotations := res.Responses[0].FaceAnnotations; len(annotations) > 0 {
		for _, a := range annotations {
			joy := a.JoyLikelihood
			anger := a.AngerLikelihood
			sorrow := a.SorrowLikelihood
			surprise := a.SurpriseLikelihood
			blurred := a.BlurredLikelihood
			score := strconv.Itoa(int(a.DetectionConfidence * 100))

			face_annotations = append(face_annotations, "{`Joy`: "+joy+", `Anger`: "+anger+", `Sorrow`: "+sorrow+", `Surprise`: "+surprise+", `Blurred`: "+blurred+", `Score`: "+score+"}")
		}
	}

	log.Println("face", strings.Join(face_annotations[:], "."))
	return strings.Join(face_annotations[:], ".")
}

//landmark-detection

func LandmarkOcr(file_dir string) string {
	var landmark_annotations []string
	res := DefaultDialer(file_dir, "LANDMARK_DETECTION")

	if annotations := res.Responses[0].LandmarkAnnotations; len(annotations) > 0 {
		for _, a := range annotations {
			landmark := a.Description
			score := strconv.Itoa(int(a.Score * 100))

			landmark_annotations = append(landmark_annotations, "{`Landmark`: "+landmark+", `Score`: "+score+"}")
		}
	}

	log.Println("landmark", strings.Join(landmark_annotations[:], "."))
	return strings.Join(landmark_annotations[:], ".")
}

//text-detection

func TextOcr(file_dir string) string {
	var text_annotations []string
	res := DefaultDialer(file_dir, "TEXT_DETECTION")

	if annotations := res.Responses[0].TextAnnotations; len(annotations) > 0 {
		for _, a := range annotations {
			text := a.Description

			text_annotations = append(text_annotations, "{`Text`: "+text+"}")
		}
	}

	log.Println("text", strings.Join(text_annotations[:], "."))
	return strings.Join(text_annotations[:], ".")
}

//safe-search-detection

func SafeSearchOcr(file_dir string) (string, string, string) {
	res := DefaultDialer(file_dir, "SAFE_SEARCH_DETECTION")

	a := res.Responses[0].SafeSearchAnnotation
	adult := a.Adult
	medical := a.Medical
	racy := a.Racy
	spoof := a.Spoof
	violence := a.Violence

	log.Println("safety", "{`Adult`: "+adult+", `Medical`: "+medical+", `Racy`: "+racy+", `Spoof`: "+spoof+", `Violence`: "+violence+"}")
	return "{`Adult`: " + adult + ", `Medical`: " + medical + ", `Racy`: " + racy + ", `Spoof`: " + spoof + ", `Violence`: " + violence + "}", adult, violence
}

func ExplicitVideoContent(file_dir string) bool {
	var IsExplicit bool = false

	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		log.Println(err)
	}

	defer client.Close()

	file, err := ioutil.ReadFile(file_dir)
	if err != nil {
		log.Println(err)
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputContent: file,
		Features:     []videopb.Feature{videopb.Feature_EXPLICIT_CONTENT_DETECTION},
	})
	if err != nil {
		log.Println(err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		log.Println(err)
	}

	result := resp.AnnotationResults[0].ExplicitAnnotation

	for _, frame := range result.Frames {
		if frame.PornographyLikelihood.String() == "VERY_LIKELY" || frame.PornographyLikelihood.String() == "LIKELY" {
			IsExplicit = true
		}
	}

	return IsExplicit
}
