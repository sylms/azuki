package usecase

import "github.com/sylms/azuki/domain"

type CourseUseCase interface {
	Search(domain.CourseQuery) ([]*domain.Course, error)
	Facet(domain.CourseQuery) ([]*domain.Facet, error)
	Update(domain.UpdateJSON) error
}

type courseUseCase struct {
	repo domain.CourseRepository
}

func NewCourseUseCase(repo domain.CourseRepository) CourseUseCase {
	return &courseUseCase{
		repo: repo,
	}
}

func (uc *courseUseCase) Search(query domain.CourseQuery) ([]*domain.Course, error) {
	courses, err := uc.repo.Search(query)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (uc *courseUseCase) Facet(query domain.CourseQuery) ([]*domain.Facet, error) {
	facets, err := uc.repo.Facet(query)
	if err != nil {
		return nil, err
	}
	return facets, nil
}

func (uc *courseUseCase) Update(query domain.UpdateJSON) error {
	err := uc.repo.Update(query)
	if err != nil {
		return err
	}
	return nil
}
