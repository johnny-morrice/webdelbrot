package nativesharedregion

import (
    "functorama.com/demo/sharedregion"
    "functorama.com/demo/nativesequence"
    "functorama.com/demo/nativeregion"
)

type NativeSharedRegion struct {
    nativeregion.NativeRegionProxy
    threadId uint
    prototypes []nativeregion.NativeRegionNumerics
    sequencePrototypes []nativesequence.NativeSequenceNumerics
}

func CreateNativeSharedRegion(numerics *nativeregion.NativeRegionNumerics, jobs uint) NativeSharedRegion {
    shared := NativeSharedRegion{
        prototypes: make([]nativeregion.NativeRegionNumerics, jobs),
        sequencePrototypes: make([]nativesequence.NativeSequenceNumerics, jobs),
    }
    for i := uint(0); i < jobs; i++ {
        shared.prototypes[i] = *numerics
        shared.sequencePrototypes[i] = *shared.prototypes[i].SequenceNumerics
    }
    initLocal := &shared.prototypes[0]
    shared.NativeRegionProxy = nativeregion.NativeRegionProxy{
        NativeRegionNumerics: initLocal,
        Region: initLocal.Region,
    }

    return shared
}

func (shared NativeSharedRegion) GrabThreadPrototype(threadId uint) {
    shared.NativeRegionProxy.NativeRegionNumerics = &shared.prototypes[threadId]
    shared.threadId = threadId
}

func (shared NativeSharedRegion) SharedChildren() []sharedregion.SharedRegionNumerics {
    localRegions := shared.NativeChildRegions()
    sharedChildren := make([]sharedregion.SharedRegionNumerics, len(localRegions))
    myCore := shared.NativeRegionProxy.NativeRegionNumerics
    for i, child := range localRegions {
        sharedChildren[i] = NativeSharedRegion{
            NativeRegionProxy: nativeregion.NativeRegionProxy{
                Region: child,
                NativeRegionNumerics: myCore,
            },
            prototypes: shared.prototypes,
        }
    }
    return sharedChildren
}

func (shared NativeSharedRegion) SharedRegionSequence() sharedregion.SharedSequenceNumerics {
    return NativeSharedSequence{
        nativeregion.NativeSequenceProxy{
            Region: shared.NativeRegionProxy.Region,
            NativeSequenceNumerics: &shared.sequencePrototypes[shared.threadId],
        },
        shared.sequencePrototypes,
    }
}

type NativeSharedSequence struct {
    nativeregion.NativeSequenceProxy
    prototypes []nativesequence.NativeSequenceNumerics
}

func (shared NativeSharedSequence) GrabThreadPrototype(threadId uint) {
    shared.NativeSequenceProxy.NativeSequenceNumerics = &shared.prototypes[threadId]
}