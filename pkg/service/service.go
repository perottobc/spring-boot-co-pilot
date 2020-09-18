package service

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

type Context struct {
	Recursive       bool
	Overwrite       bool
	DryRun          bool
	TargetDirectory string
	PomPairs        []maven.PomPair
	Err             error
}

var log = logger.Context()

func Write(overwrite bool, pair maven.PomPair) error {
	var writeToFile = pair.PomFile
	if !overwrite {
		writeToFile = pair.PomFile + ".new"
	}
	if err := maven.SortAndWritePom(pair.Model, writeToFile); err != nil {
		return err
	}

	return nil
}

func PomFileToTargetDirectory(pomFile string) string {
	pomFilePathParts := strings.Split(pomFile, "/")
	return strings.Join(pomFilePathParts[:len(pomFilePathParts)-1], "/")
}

func (ctx *Context) FindAndPopulatePomModels() {
	excludes := []string{
		"flattened-pom.xml",
		"/target/",
	}

	if ctx.Recursive {
		pomFiles, err := file.FindAll("pom.xml", excludes, ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}
		for _, pomFile := range pomFiles {
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}
			ctx.PomPairs = append(ctx.PomPairs, maven.PomPair{
				Model:   model,
				PomFile: pomFile,
			})
		}
	} else {
		pomFile := fmt.Sprintf("%s/pom.xml", ctx.TargetDirectory)
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Warnln(err)
			return
		}
		ctx.PomPairs = append(ctx.PomPairs, maven.PomPair{
			Model:   model,
			PomFile: pomFile,
		})
	}
}

func (ctx Context) OnEachPomProject(description string, do func(pair maven.PomPair, args ...interface{}) error) {
	if ctx.PomPairs == nil {
		log.Errorln("could not find any pom models in the context")
		return
	}

	for _, pair := range ctx.PomPairs {
		log.Info(logger.White(fmt.Sprintf("%s for pom file %s", description, pair.PomFile)))

		if do != nil {
			err := do(pair)
			if err != nil {
				log.Warnln(err)
				continue
			}
		}

		if !ctx.DryRun {
			if err := Write(ctx.Overwrite, pair); err != nil {
				log.Warnln(err)
			}
		}
	}
}
