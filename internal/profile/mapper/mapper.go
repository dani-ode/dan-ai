// internal/profile/mapper/mapper.go
package mapper

import (
	"dan-ai/internal/profile/entity"
	pb "dan-ai/proto/profile"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProto maps a domain Profile entity to a protobuf Profile message.
func ToProto(e *entity.Profile) *pb.Profile {
	if e == nil {
		return nil
	}
	return &pb.Profile{
		Id:           e.ID,
		FullName:     e.FullName,
		Headline:     e.Headline,
		Bio:          e.Bio,
		Email:        e.Email,
		Phone:        e.Phone,
		Location:     e.Location,
		Github:       e.Github,
		Linkedin:     e.Linkedin,
		Website:      e.Website,
		Avatar:       e.Avatar,
		ResumeUrl:    e.ResumeURL,
		Availability: e.Availability,
		Timezone:     e.Timezone,
		CreatedAt:    timestamppb.New(e.CreatedAt),
		UpdatedAt:    timestamppb.New(e.UpdatedAt),
	}
}

// ToEntityFromCreate maps a CreateProfileRequest protobuf message to a domain Profile entity.
func ToEntityFromCreate(r *pb.CreateProfileRequest) *entity.Profile {
	if r == nil {
		return nil
	}
	return &entity.Profile{
		FullName:     r.FullName,
		Headline:     r.Headline,
		Bio:          r.Bio,
		Email:        r.Email,
		Phone:        r.Phone,
		Location:     r.Location,
		Github:       r.Github,
		Linkedin:     r.Linkedin,
		Website:      r.Website,
		Avatar:       r.Avatar,
		ResumeURL:    r.ResumeUrl,
		Availability: r.Availability,
		Timezone:     r.Timezone,
	}
}

// ToEntityFromUpdate maps an UpdateProfileRequest protobuf message to a domain Profile entity.
func ToEntityFromUpdate(r *pb.UpdateProfileRequest) *entity.Profile {
	if r == nil {
		return nil
	}
	return &entity.Profile{
		ID:           r.Id,
		FullName:     r.FullName,
		Headline:     r.Headline,
		Bio:          r.Bio,
		Email:        r.Email,
		Phone:        r.Phone,
		Location:     r.Location,
		Github:       r.Github,
		Linkedin:     r.Linkedin,
		Website:      r.Website,
		Avatar:       r.Avatar,
		ResumeURL:    r.ResumeUrl,
		Availability: r.Availability,
		Timezone:     r.Timezone,
	}
}
